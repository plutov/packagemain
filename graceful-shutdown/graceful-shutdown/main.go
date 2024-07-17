package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis"
)

var wg sync.WaitGroup

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// not necessary to manually cancel with RegisterOnShutdown(cancel) below
	// defer cancel()
	signalCtx, signalCtxStop := signal.NotifyContext(ctx,
		os.Interrupt,    // interrupt = SIGINT = Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)
	defer signalCtxStop()

	redisdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	// inherit from signal cancellation context
	server := http.Server{
		Addr:        ":8080",
		BaseContext: func(_ net.Listener) context.Context { return signalCtx },
	}

	// server.RegisterOnShutdown registers a function to call on Shutdown.
	// This can be used to gracefully shutdown connections that have
	// undergone ALPN protocol upgrade (HTTP2 or HTTP3)
	// or that have been hijacked (websockets).
	//
	// The RegisterOnShutdown's argument function should start
	// protocol-specific graceful shutdown,
	// but should not wait for shutdown to complete.
	//
	// Also needed if your BaseContext is more complex you might want
	// to use this instead of calling cancel manually
	server.RegisterOnShutdown(cancel)

	http.HandleFunc("/incr", func(w http.ResponseWriter, r *http.Request) {
		//TODO: No longer necessary
		wg.Add(1)
		go processRequest(redisdb)
		w.WriteHeader(http.StatusOK)
	})

	// make it a goroutine
	go server.ListenAndServe()

	// listen for the interrupt signal.
	<-signalCtx.Done()

	// Give outstanding requests a deadline for completion.
	const timeout = 5 * time.Second
	shutdownCtx, shutdownRelease := context.WithTimeout(signalCtx, timeout)
	defer shutdownRelease()

	// signalCtxStop the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("could not shutdown: %v\n", err)
	}

	//TODO: No longer necessary
	// wait for all goroutines to finish
	wg.Wait()

	// close redis connection
	redisdb.Close()

	os.Exit(0)
}

func processRequest(redisdb *redis.Client) {
	//TODO: No longer necessary
	defer wg.Done()

	// simulate some business logic here
	time.Sleep(time.Second * 5)
	redisdb.Incr("counter")
}
