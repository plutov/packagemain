package main

import (
	"log"
	"math/rand"
	"net/http"
	neturl "net/url"
)

type server struct {
	DB    DB
	Cache Cache
}

const (
	charset   = "abcdefghijklmnopqrstuvwxyz01234567890"
	keyLength = 8
)

func generateKey() string {
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func NewServer(db DB, cache Cache) (*server, error) {
	if err := db.Init(); err != nil {
		return nil, err
	}

	if err := cache.Init(); err != nil {
		return nil, err
	}

	return &server{DB: db, Cache: cache}, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/create":
		url := req.URL.Query().Get("url")
		if _, err := neturl.ParseRequestURI(url); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		key := generateKey()
		if err := s.DB.StoreURL(url, key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(key))
	case "/get":
		key := req.URL.Query().Get("key")
		if url, err := s.Cache.Get(key); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(url))
		}

		url, err := s.DB.GetURL(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s.Cache.Set(key, url)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(url))
	}
}

func main() {
	s, err := NewServer(&MongoDB{}, &Redis{})
	if err != nil {
		log.Fatalf("error initializing server: %v", err)
	}

	http.ListenAndServe(":8080", s)
}
