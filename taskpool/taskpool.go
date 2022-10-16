package taskpool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"taskpool/option"
	"time"
)

type TaskPool interface {
	Start() error
	Submit(ctx context.Context, task Task) error
	Shutdown() (<-chan struct{}, error)
	ShutdownNow() ([]Task, error)
}
type Task interface {
	Run(ctx context.Context) error
}

var (
	stateCreated int32 = 1
	stateRunning int32 = 2
	stateClosing int32 = 3
	stateStopped int32 = 4
	stateLocked  int32 = 5

	errTaskPoolIsNotRunning = errors.New("ekit: TaskPool未运行")
	errTaskPoolIsClosing    = errors.New("ekit：TaskPool关闭中")
	errTaskPoolIsStopped    = errors.New("ekit: TaskPool已停止")
	errTaskPoolIsStarted    = errors.New("ekit：TaskPool已运行")
	errTaskIsInvalid        = errors.New("ekit: Task非法")
	errTaskRunningPanic     = errors.New("ekit: Task运行时异常")

	errInvalidArgument = errors.New("ekit: 参数非法")

	_            TaskPool = &OnDemandBlockTaskPool{}
	panicBuffLen          = 2048

	defaultMaxIdleTime = 10 * time.Second
)

type TaskFunc func(ctx context.Context) error

func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

type taskWrapper struct {
	t Task
}

func (tw *taskWrapper) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, panicBuffLen)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("%w: %s", errTaskRunningPanic, fmt.Sprintf("[PANIC]:\t%+v\n%s\n", r, buf))
		}
	}()
	return tw.t.Run(ctx)
}

type group struct {
	mp map[int]int
	n  int32
	mu sync.RWMutex
}

func (g *group) isIn(id int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.mp[id]
	return ok
}

func (g *group) add(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.mp[id]; !ok {
		g.mp[id] = 1
		g.n++
	}
}

func (g *group) delete(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.mp[id]; ok {
		g.n--
	}
	delete(g.mp, id)
}
func (g *group) size() int32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.n
}

type OnDemandBlockTaskPool struct {
	state             int32
	queue             chan Task
	numGoRunningTasks int32
	totalGo           int32
	mutex             sync.RWMutex
	initGo            int32
	coreGo            int32
	maxGo             int32
	queueBacklogRate  float64
	shutdownOnce      sync.Once
	id                int32
	shutdownDone      chan struct{}
	shutdownNowCtx    context.Context
	shutdownCancel    context.CancelFunc
	maxIdleTime       time.Duration
	timeoutGroup      *group
}

func NewOnDemandBlockTaskPool(initGo int, queueSize int, opts ...option.Option[OnDemandBlockTaskPool]) (*OnDemandBlockTaskPool, error) {
	if initGo < 1 {
		return nil, fmt.Errorf("%w：initGo应该大于0", errInvalidArgument)
	}
	if queueSize < 0 {
		return nil, fmt.Errorf("%w：queueSize应该大于等于0", errInvalidArgument)
	}
	b := &OnDemandBlockTaskPool{
		queue:        make(chan Task, queueSize),
		shutdownDone: make(chan struct{}, 1),
		initGo:       int32(initGo),
		coreGo:       int32(initGo),
		maxGo:        int32(initGo),
		maxIdleTime:  defaultMaxIdleTime,
	}
	b.shutdownNowCtx, b.shutdownCancel = context.WithCancel(context.Background())
	atomic.StoreInt32(&b.state, stateCreated)
	option.Apply(b, opts...)
	if b.coreGo != b.initGo && b.maxGo == b.initGo {
		b.maxGo = b.coreGo
	} else if b.coreGo == b.initGo && b.maxGo != b.initGo {
		b.coreGo = b.maxGo
	}
	if !(b.initGo <= b.coreGo && b.coreGo <= b.maxGo) {
		return nil, fmt.Errorf("%w: 需要满足initGo <= coreGo <= maxGo条件 ", errInvalidArgument)
	}
	b.timeoutGroup = &group{mp: make(map[int]int)}
	if b.queueBacklogRate < float64(0) || float64(1) < b.queueBacklogRate {
		return nil, fmt.Errorf("%w,queueBacklogRate合法范围[0,1.0]", errInvalidArgument)
	}
	return b, nil
}

func (b *OnDemandBlockTaskPool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return fmt.Errorf("%w", errTaskIsInvalid)
	}
	for {
		if atomic.LoadInt32(&b.state) == stateClosing {
			return errTaskPoolIsClosing
		}
		if atomic.LoadInt32(&b.state) == stateStopped {
			return errTaskPoolIsStopped
		}
		task = &taskWrapper{t: task}

	}
}
func (b *OnDemandBlockTaskPool) allowToCreateGoroutine() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	if b.totalGo == b.maxGo {
		return false
	}
	rate := float64(len(b.queue)) / float64(cap(b.queue))
	if rate == 0 || rate < b.queueBacklogRate {
		log.Println("rate == 0", rate == 0, "rate", rate, "<", b.queueBacklogRate)
		return false
	}
	return true
}

func (b *OnDemandBlockTaskPool) trySubmit(ctx context.Context, task Task, state int32) (bool, error) {
	if atomic.CompareAndSwapInt32(&b.state, state, stateLocked) {
		defer atomic.CompareAndSwapInt32(&b.state, stateLocked, state)
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("%w", ctx.Err())
		case b.queue <- task:
			if state == stateRunning && b.allowToCreateGoroutine() {
				b.increaseTotalGo(1)
				go b.goroutine(int(atomic.LoadInt32(&b.id)))
				newId := atomic.AddInt32(&b.id, 1)
				log.Println("create go", newId-1)
			}
			return true, nil
		default:
			return false, nil
		}
	}

	return false, nil
}

func (b *OnDemandBlockTaskPool) increaseTotalGo(n int32) {
	b.mutex.Lock()
	b.totalGo += n
	b.mutex.Unlock()
}
func (b *OnDemandBlockTaskPool) decreaseTotalGo(n int32) {
	b.mutex.Lock()
	b.totalGo -= n
	b.mutex.Unlock()
}
func (b *OnDemandBlockTaskPool) goroutine(id int) {
	idleTimer := time.NewTimer(0)
	if !idleTimer.Stop() {
		<-idleTimer.C
	}
	for {
		select {
		case <-b.shutdownNowCtx.Done():
			b.decreaseTotalGo(1) // ?为什么不用atomic
			return
		case <-idleTimer.C:
			b.mutex.Lock()
			b.totalGo--
			b.timeoutGroup.delete(id)
			b.mutex.Unlock()
			return
		case task, ok := <-b.queue:
			if b.timeoutGroup.isIn(id) {
				b.timeoutGroup.delete(id)
				if !idleTimer.Stop() {
					<-idleTimer.C
				}
			}
			atomic.AddInt32(&b.numGoRunningTasks, 1)
			if !ok {
				if atomic.CompareAndSwapInt32(&b.numGoRunningTasks, 1, 0) && atomic.LoadInt32(&b.state) == stateClosing {
					b.shutdownOnce.Do(func() {
						atomic.CompareAndSwapInt32(&b.state, stateClosing, stateStopped)
						b.shutdownDone <- struct{}{}
						close(b.shutdownDone)
					})
					b.decreaseTotalGo(1)
					return
				}
				atomic.AddInt32(&b.numGoRunningTasks, -1)
				b.decreaseTotalGo(1)
				return
			}
			_ = task.Run(b.shutdownNowCtx)
			atomic.AddInt32(&b.numGoRunningTasks, -1)
			b.mutex.Lock()
			if b.coreGo < b.totalGo && (len(b.queue) == 0 || int32(len(b.queue)) < b.totalGo) {
				b.totalGo--
				b.mutex.Unlock()
				return
			}
			if b.initGo < b.totalGo-b.timeoutGroup.size() {
				idleTimer = time.NewTimer(b.maxIdleTime)
				b.timeoutGroup.add(id)
			}
			b.mutex.Unlock()
		}
	}
}
