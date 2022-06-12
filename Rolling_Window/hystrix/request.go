package hystrix

import (
	"sync"
	"time"
)

type Requests struct {
	lock sync.RWMutex
	Total int
	Failed int
	Timestamp time.Time
}

func NewResquests() *Requests {
	return &Requests{
		Timestamp: time.Now(),
	}
}


func (r *Requests)Record(result bool){
	r.lock.Lock()
	defer r.lock.Unlock()
	if !result {
		r.Failed++
	}
	r.Total++

}