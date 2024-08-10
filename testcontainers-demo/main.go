package main

import "net/http"

type server struct{}

func (r *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	http.ListenAndServe(":8080", &server{})
}
