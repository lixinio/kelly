package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lixinio/kelly"
)

func waitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}

func watchSignal(cancel context.CancelFunc) {
	stop := make(chan struct{})
	waitSignal(stop)

	<-stop
	fmt.Println("Signal caught; shutting down server")
	cancel()
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	r := func(c *kelly.Context) {
		c.WriteIndentedJSON(http.StatusOK, kelly.H{
			"code": "0",
		})
	}

	router1 := kelly.New(nil)
	router1.GET("/", r)
	router2 := kelly.New(nil)
	router2.GET("/", r)

	go watchSignal(cancel)

	go func() {
		defer wg.Done()
		fmt.Println("router 1 started")
		router1.RunContext(ctx, ":9998")
		fmt.Println("router 1 finished")
	}()
	go func() {
		defer wg.Done()
		fmt.Println("router 2 started")
		router2.RunContext(ctx, ":9999")
		fmt.Println("router 2 finished")
	}()

	fmt.Println("main waiting")
	wg.Wait()
	fmt.Println("main finished")
}
