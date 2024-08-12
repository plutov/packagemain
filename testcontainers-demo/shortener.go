package main

import (
	"fmt"
	"math/rand"
	neturl "net/url"
)

func StoreURL(db DB, url string) (string, error) {
	if _, err := neturl.ParseRequestURI(url); err != nil {
		return "", err
	}

	key := generateKey()

	if err := db.StoreURL(url, key); err != nil {
		return "", err
	}

	return key, nil
}

func GetURL(db DB, cache Cache, key string) (string, error) {
	if len(key) != keyLength {
		return "", fmt.Errorf("invalid key length, expecting: %d", keyLength)
	}

	if url, ok := cache.Get(key); ok {
		return url, nil
	}

	url, err := db.GetURL(key)
	if err != nil {
		return "", err
	}

	cache.Set(key, url)

	return url, nil
}

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 8
)

func generateKey() string {
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
