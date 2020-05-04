package testable

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Repo struct {
	StargazersCount int `json:"stargazers_count"`
}

type RepositoriesAPI interface {
	GetRepos(username string) ([]Repo, error)
}

type Mock struct{}

func (m *Mock) GetRepos(username string) ([]Repo, error) {
	return []Repo{
		Repo{
			StargazersCount: 2,
		},
		Repo{
			StargazersCount: 6,
		},
	}, nil
}

type GitHub struct{}

func (g *GitHub) GetRepos(username string) ([]Repo, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/repos", username))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	repos := []Repo{}
	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}
