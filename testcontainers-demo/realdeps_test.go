//go:build realdeps
// +build realdeps

package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerWithRealDependencies(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("REDIS_URI", "redis://localhost:6379")

	s, err := NewServer(&MongoDB{}, &Redis{})
	assert.NoError(t, err)

	srv := httptest.NewServer(s)
	defer srv.Close()

	testServer(srv, t)
}
