package app

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
