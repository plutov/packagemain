//go:build integration
// +build integration

package main

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestServerWithTestcontainers(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "docker.io/mongodb/mongodb-community-server:7.0-ubi8")
	assert.NoError(t, err)
	defer mongodbContainer.Terminate(ctx)

	redisContainer, err := redis.Run(ctx, "docker.io/redis:7.4-alpine")
	assert.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	mongodbEndpoint, _ := mongodbContainer.Endpoint(ctx, "")
	redisEndpoint, _ := redisContainer.Endpoint(ctx, "")

	os.Setenv("MONGO_URI", "mongodb://"+mongodbEndpoint)
	os.Setenv("REDIS_URI", "redis://"+redisEndpoint)

	s, err := NewServer(&MongoDB{}, &Redis{})
	assert.NoError(t, err)

	srv := httptest.NewServer(s)
	defer srv.Close()

	testServer(srv, t)
}
