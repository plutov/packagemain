package main

import (
	"net/http"
)

type server struct {
	DB    DB
	Cache Cache
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

		w.Write([]byte(key))
	case "/get":
		key := req.URL.Query().Get("key")
		url, err := GetURL(r.DB, r.Cache, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(url))
	}
}

func main() {
	http.ListenAndServe(":8080", &server{})
}
