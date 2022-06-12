package hystrix

import (
	"fmt"
	"sync"
	"time"
)

type RollingWindow struct {
	lock            sync.RWMutex
	broken          bool
	size            int
	Requestsli      []*Requests
	reqThreshold    int
	failedThreshold float64   //触发熔断的阈值
	lastBreakTime   time.Time //上次熔断发生的时间
	seeker          bool
	brokenTimeGap   time.Duration //恢复时间
}
func NewRollingWindow(size int,reqThreshold int,failedThreshold float64,brokenTimeGap time.Duration)*RollingWindow{
	return &RollingWindow{
		size: size,
		reqThreshold: reqThreshold,
		failedThreshold: failedThreshold,
		brokenTimeGap: brokenTimeGap,
		Requestsli: make([]*Requests,0,size),
	}
}

//添加一个新请求结果
func (R *RollingWindow)ADDNewRequests() {
	R.lock.Lock()
	defer R.lock.Unlock()
	if len(R.Requestsli) >= R.size {
		R.Requestsli = R.Requestsli[1:]
	}
	R.Requestsli = append(R.Requestsli,NewResquests())
}
//获取最新的请求结果
func (R *RollingWindow)GetRquests()*Requests{

	return R.Requestsli[len(R.Requestsli)-1]
}
//记录结果
func (R *RollingWindow)RecordReqResult(result bool){
	R.GetRquests().Record(result)
}

//展示当前滑动窗口所有请求集的结果
func (R *RollingWindow)ShowRequestsli()  {
	for _,v := range R.Requestsli {
		fmt.Printf("time:%v  ,total:%d ,failed: %d\n",v.Timestamp,v.Total,v.Failed)
	}
}
func (R *RollingWindow)Start() {
	go func() {
		for {
			R.ADDNewRequests()
			time.Sleep(time.Millisecond*100)
		}
	}()
}
//根据当前窗口的判断是否需要进入熔断

func (R *RollingWindow)Judgment_blown()bool{
	R.lock.RLock()
	defer R.lock.RUnlock()
	failed,total := 0,0

	for _,v := range R.Requestsli {
		failed += v.Failed
		total += v.Total
	}
	if total >= R.reqThreshold {

		return true
	}
	if float64(failed)/float64(total) >= R.failedThreshold{
		return true
	}
	return false
}
//监控滑动窗口的总失败次数和是否开启熔断
func (R *RollingWindow)Monitor(){
	go func(){
		for {
		if R.broken {
			if R.JudgeOverbroken() {
			R.lock.Lock()
			R.broken = false
			R.lock.Unlock()
			}
		}
		if R.Judgment_blown(){
			R.lock.Lock()
			R.broken = true
			R.lastBreakTime = time.Now()
			R.lock.Unlock()
		}
		}
	}()
}
func (R *RollingWindow)JudgeOverbroken()bool{
	return time.Since(R.lastBreakTime) > R.brokenTimeGap
}

func (R *RollingWindow)Broken()bool{
	return R.broken
}