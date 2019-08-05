package main

import (
	"log"
	"net/http"
)

var limiter = NewIPRateLimiter(1, 5)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	if err := http.ListenAndServe(":8888", limitMiddleware(mux)); err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
}

func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetVisitor(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	// Some very expensive database call
	w.Write([]byte("alles gut"))
}
