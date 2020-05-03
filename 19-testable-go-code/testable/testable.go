package testable

func StrInSlice(slice []string, find string) bool {
	for _, v := range slice {
		if v == find {
			return true
		}
	}

	return false
}

func GetAverageStarsPerRepo(repositoriesAPI RepositoriesAPI, username string) (float64, error) {
	repos, err := repositoriesAPI.GetRepos(username)
	if err != nil {
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
