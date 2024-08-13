package main

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Init() error
	Set(key string, val string) error
	Get(key string) (string, bool)
}

type Redis struct {
	client *redis.Client
}

func (r *Redis) Init() error {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		return err
	}

	r.client = redis.NewClient(opts)
	return nil
}

func (r *Redis) Set(key string, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return r.client.Set(ctx, key, val, 0).Err()
}

func (r *Redis) Get(key string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}

	return data, true
}
