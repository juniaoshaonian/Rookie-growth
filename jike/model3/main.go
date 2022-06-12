package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func  main(){
	g, ctx := errgroup.WithContext(context.Background())

	h:= http.NewServeMux()
	h.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("hi"))
	})
	shutdown := make(chan struct{})
	h.HandleFunc("/stop",func(w http.ResponseWriter,r *http.Request){
		shutdown <- struct{}{}
	})
	server := http.Server{
		Handler: h,
		Addr: ":8080",
	}
	g.Go(func() error {
		return server.ListenAndServe()
	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("errgroup exit...")
		case <-shutdown:
			fmt.Println("server will out...")
		}

		Ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		fmt.Println("shutting down server...")
		return server.Shutdown(Ctx)
	})
	g.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.New(fmt.Sprintf("os exited%v",sig))
		}
	})

	fmt.Printf("errgroup exiting: %+v\n", g.Wait())
}
