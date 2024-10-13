package main

import (
	"context"
	"log"
	"net/http"
	"time"

	gogithub "github.com/google/go-github/v65/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/redirect", redirect)
	mux.HandleFunc("/callback", callback)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("unable to start the server: %w", err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, getRedirectURL(), http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// TODO: validate state
	// state := r.URL.Query().Get("state")

	conf := getOAuthConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("unable to get token: %w", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repos, err := getCurrentUserRepos(token.AccessToken)
	if err != nil {
		log.Printf("unable to get repos: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, r := range repos {
		w.Write([]byte(r.GetFullName() + ", "))
	}
	w.WriteHeader(http.StatusOK)
}

func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		// TODO: use env vars instead
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user:email", "repo:public_repo"},
		Endpoint:     github.Endpoint,
	}
}

func getRedirectURL() string {
	conf := getOAuthConfig()
	return conf.AuthCodeURL("state")
}

func getCurrentUserRepos(accessToken string) ([]*gogithub.Repository, error) {
	client := gogithub.NewClient(nil).WithAuthToken(accessToken)

	opt := &gogithub.RepositoryListByAuthenticatedUserOptions{
		Affiliation: "owner",
	}
	opt.PerPage = 50
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)

	return repos, err
}
