package main

import (
	"log"
	"net/http"
)

type server struct {
	DB    DB
	Cache Cache
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

func (r *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/create":
		url := req.URL.Query().Get("url")
		key, err := StoreURL(r.DB, url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(key))
	case "/get":
		key := req.URL.Query().Get("key")
		url, err := GetURL(r.DB, r.Cache, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

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
