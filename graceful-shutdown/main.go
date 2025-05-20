package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	cache := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	http.HandleFunc("/incr", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // simulate some delay
		cache.Incr("counter")
		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{
		Addr: ":8080",
	}

	server.ListenAndServe()
}
