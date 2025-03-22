package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Init() error
	Set(key string, val string) error
	Get(key string) (string, error)
}

type Redis struct {
	client *redis.Client
}

func (r *Redis) Init() error {
	opts, err := redis.ParseURL(os.Getenv("REDIS_ADDR"))
	if err != nil {
		return err
	}

	r.client = redis.NewClient(opts)
	return nil
}

func (r *Redis) Set(key string, val string) error {
	return r.client.Set(context.Background(), key, val, 0).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}
