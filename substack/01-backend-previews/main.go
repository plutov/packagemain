package main

import (
	"encoding/json"
	"net/http"
)

type mainHandler struct{}
type substack struct {
	Name    string   `json:"name"`
	Authors []string `json:"authors"`
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func (h *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sub := substack{
		Name: "package main",
		Authors: []string{
			"Alex Pliutau",
			"Julien Singler",
		},
	}
	jsonBytes, err := json.Marshal(sub)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &mainHandler{})
	http.ListenAndServe(":8080", mux)
}
