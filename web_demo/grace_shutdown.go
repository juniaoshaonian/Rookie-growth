package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type gracefulshutdown struct {
	reqcnt int64
	closing int32
	reqchan chan struct{}
}
var ErrorHookTimeout error = errors.New("call hook timeout")
func Newgracefulshutdown()*gracefulshutdown{
	return &gracefulshutdown{
		reqchan: make(chan struct{},1),
	}
}
func (g *gracefulshutdown)Shutdownfilter(n filter)filter{
	return func(c *Context){
		cl := atomic.LoadInt32(&g.closing)
		if cl > 0 {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		atomic.AddInt64(&g.reqcnt,1)
		n(c)
		atomic.AddInt64(&g.reqcnt,-1)
		if cl > 0 && g.reqcnt == 0{
			g.reqchan <- struct{}{}
		}
	}
}
func (g *gracefulshutdown)RejectNewRequestAndWaiting(ctx context.Context)error{
	atomic.AddInt32(&g.closing,1)
	if atomic.LoadInt64(&g.reqcnt) == 0 {
		return nil
	}
	done := ctx.Done()
	select{
	case <-done:
		fmt.Println("超时了，还没等到所有请求")
		return ErrorHookTimeout
	case <- g.reqchan:
		fmt.Println("全部处理完了")
	}
	return nil
}
func WaitForShutdown(hooks... hook){
	signals := make(chan os.Signal,1)
	signal.Notify(signals,ShutdownSignals...)
	select {
	case sig := <-signals:
		fmt.Printf("get signal #{sig},application will shutdown\n" )
		time.AfterFunc(time.Minute*10,func(){
			fmt.Printf("Shutdown gracefully timeout, application will shutdown immediately. ")
		})
		for _,h := range hooks{
			ctx,cancel := context.WithTimeout(context.Background(),time.Second*30)
			err := h(ctx)
			if err != nil {
				fmt.Println("failed to run hook ")
			}
			cancel()
		}
		os.Exit(0)
	}


}