package main

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/plutov/packagemain/urlshortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServerWithMocks(t *testing.T) {
	mockDB := mocks.NewDB(t)
	mockCache := mocks.NewCache(t)

	mockDB.EXPECT().Init().Return(nil)
	mockCache.EXPECT().Init().Return(nil)

	mockDB.EXPECT().StoreURL(mock.Anything, mock.Anything).Return(nil)

	mockCache.EXPECT().Get("invalidkey").Return("", errors.New("not found"))
	mockDB.EXPECT().GetURL("invalidkey").Return("", errors.New("not found"))

	mockCache.EXPECT().Get(mock.Anything).Return("", errors.New("not found"))
	mockDB.EXPECT().GetURL(mock.Anything).Return("https://packagemain.tech", nil)
	mockCache.EXPECT().Set(mock.Anything, mock.Anything).Return(nil)

	s, err := NewServer(mockDB, mockCache)
	assert.NoError(t, err)

	srv := httptest.NewServer(s)
	defer srv.Close()

	testServer(srv, t)
}
