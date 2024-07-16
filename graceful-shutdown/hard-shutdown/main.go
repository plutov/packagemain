package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	redisdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	redisdb.Set("counter", "0", 0)

	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/incr", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2) // Simulate some business logic here
		redisdb.Incr("counter")
		w.WriteHeader(http.StatusOK)
	})

	server.ListenAndServe()
}
