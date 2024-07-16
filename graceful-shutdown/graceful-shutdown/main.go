package main

import (
	"context"
	"log"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	redisdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/incr", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		go processRequest(redisdb)
		w.WriteHeader(http.StatusOK)
	})

	// make it a goroutine
	go server.ListenAndServe()

	// listen for the interrupt signal.
	<-ctx.Done()

	// stop the server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("could not shutdown: %v\n", err)
	}

	// wait for all goroutines to finish
	wg.Wait()

	// close redis connection
	redisdb.Close()

	os.Exit(0)
}

func processRequest(redisdb *redis.Client) {
	defer wg.Done()

	// simulate some business logic here
	time.Sleep(time.Second * 5)
	redisdb.Incr("counter")
}
