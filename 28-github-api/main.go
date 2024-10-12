package main

import (
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/redirect", redirect)
	mux.HandleFunc("/callback", callback)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
}

func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		// TODO: use env vars
		ClientID:     "Ov23liNqJR9Ysv2Q967f",
		ClientSecret: "da0cc01edf3043a3a11e806ff5dbfb9df80feb1a",
		Scopes:       []string{"user:email", "repo:public_repo"},
		Endpoint:     github.Endpoint,
	}
}

func getRedirectURL() string {
	conf := getOAuthConfig()
	return conf.AuthCodeURL("state")
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, getRedirectURL(), http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
