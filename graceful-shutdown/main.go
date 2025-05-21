package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	shutdownDelay         = 5 * time.Second
	serverShutdownTimeout = 10 * time.Second
	hardShutdownWait      = 5 * time.Second
)

var (
	shutdownInProgress atomic.Bool
	cache              *redis.Client
)

func main() {
	cache = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl+\
		syscall.SIGTERM, // the normal way to politely ask a program to terminate
	)
	defer signalCtxStop()

	// readiness check
	http.HandleFunc("/", ready)
	// main endpoint
	http.HandleFunc("/incr", incr)

	baseCtx, baseCtxStop := context.WithCancel(context.Background())
	server := http.Server{
		Addr: ":8080",
		BaseContext: func(_ net.Listener) context.Context {
			// do not pass a signalCtx here, we don't want to cancel all ongoing requests immediately
			return baseCtx
		},
	}

	// run server in a goroutine
	go server.ListenAndServe()

	// listen for the interrupt signal
	<-signalCtx.Done()

	log.Println("shutdown initiated")
	shutdownInProgress.Store(true)

	// run server as is with readiness check failing for short time
	time.Sleep(shutdownDelay)

	// give server shutdown process a deadline
	shutdownCtx, shutdownCtxStop := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer shutdownCtxStop()

	// shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("could not shutdown ongoing requests: %v\n", err)
		time.Sleep(hardShutdownWait)
	}
	// cancel the server BaseContext
	baseCtxStop()

	log.Println("shutdown complete")

	// close redis connection
	cache.Close()

	os.Exit(0)
}

func incr(w http.ResponseWriter, r *http.Request) {
	// simulate some delay
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	cache.Incr(r.Context(), "counter")
	fmt.Fprintln(w, "ok")
	w.WriteHeader(http.StatusOK)
}

func ready(w http.ResponseWriter, r *http.Request) {
	if shutdownInProgress.Load() {
		fmt.Fprintln(w, "shutdown in progress")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintln(w, "ok")
	w.WriteHeader(http.StatusOK)
}
