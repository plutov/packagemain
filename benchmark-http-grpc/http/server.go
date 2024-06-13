package http_server

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
	Code    uint64 `json:"code"`
	User    *User  `json:"user"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Response{
		Code:    201,
		Message: "ok",
		User:    &user,
	})
}

func StartHTTP1() {
	http.HandleFunc("/", CreateUser)
	if err := http.ListenAndServe(":60001", nil); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func StartHTTP2() {
	h2s := &http2.Server{}
	handler := http.HandlerFunc(CreateUser)
	server := &http.Server{
		Addr:    ":60002",
		Handler: h2c.NewHandler(handler, h2s),
	}
	if err := http2.ConfigureServer(server, h2s); err != nil {
		log.Fatalf("failed to configure http2 server: %v", err)
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
