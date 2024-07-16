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

	server := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/incr", func(w http.ResponseWriter, r *http.Request) {
		go processRequest(redisdb)
		w.WriteHeader(http.StatusOK)
	})

	server.ListenAndServe()
}

func processRequest(redisdb *redis.Client) {
	// simulate some business logic here
	time.Sleep(time.Second * 5)
	redisdb.Incr("counter")
}
