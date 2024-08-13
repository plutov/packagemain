package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestResponse(method string, url string) (string, int, error) {
	req, _ := http.NewRequest(method, url, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	return strings.Trim(string(b), "\n"), res.StatusCode, nil
}

func testServer(srv *httptest.Server, t *testing.T) {
	_, code, _ := getTestResponse("GET", srv.URL+"/create?url=invalidurl")
	assert.Equal(t, http.StatusInternalServerError, code)

	key, code, _ := getTestResponse("GET", srv.URL+"/create?url=https://packagemain.tech")
	assert.Equal(t, http.StatusCreated, code)
	assert.Len(t, key, keyLength)

	_, code, _ = getTestResponse("GET", srv.URL+"/get?key=invalidkey")
	assert.Equal(t, http.StatusNotFound, code)

	url, code, _ := getTestResponse("GET", srv.URL+"/get?key="+key)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "https://packagemain.tech", url)
}
