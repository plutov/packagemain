package main

import (
	"net/http/httptest"
	"testing"

	"github.com/plutov/packagemain/testcontainers-demo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateKey(t *testing.T) {
	key := generateKey()
	assert.Equal(t, keyLength, len(key))
}

func TestServer(t *testing.T) {
	mockDB := mocks.NewDB(t)
	mockCache := mocks.NewCache(t)

	mockDB.EXPECT().StoreURL(mock.Anything, mock.Anything).Return(nil)
	mockDB.EXPECT().GetURL(mock.Anything).Return("https://packagemain.tech", nil)

	mockCache.EXPECT().Get(mock.Anything).Return("", false)
	mockCache.EXPECT().Get(mock.Anything).Return("https://packagemain.tech", true)
	mockCache.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

	s := &server{
		DB:    mockDB,
		Cache: mockCache,
	}
	srv := httptest.NewServer(s)
	defer srv.Close()

	testServer(srv, t)
}
