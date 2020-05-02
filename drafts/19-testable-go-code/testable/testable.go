package testable

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func strInSlice(slice []string, find string) bool {
	for _, v := range slice {
		if v == find {
			return true
		}
	}

	return false
}

type Repo struct {
	StargazersCount int `json:"stargazers_count"`
}

func getAveragePerRepo(username string) (float64, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/repos", username))
	if err != nil {
		return 0, err
	}

	repos := []Repo{}
	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		return 0, err
	}

	if len(repos) == 0 {
		return 0, nil
	}

	var total int
	for _, r := range repos {
		total += r.StargazersCount
	}

	return float64(total) / float64(len(repos)), nil
}
