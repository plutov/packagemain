## packagemain #19: Writing testable Go code

When I say "testable code", what I mean is code that can be easily programmatically verified. We can say that code is testable when we don't have to change the code itself when we're adding a unit test to it. It doesn't matter if you're following test-driven development or not, testable code makes your program more flexible and maintainable, due to its modularity.

Go has robust built-in testing functionality, so in most cases you don't need to import any third-party testing packages. So start clean in the beginning, and if it's not enough, you can add helper package later ([assert](https://pkg.go.dev/github.com/stretchr/testify/assert) for example).

### SOLID

First of all, understanding of SOLID principles will help you with writing testable code. I won't go into details, but Single Responsibility and Dependency Inversion will help you a lot.

For example, it's much easier and cleaner to test small function which does only one thing. For example, function `StrInSlice` is perfectly testable function, it's determenistic, so for any given input there is only one correct output.

```go
func StrInSlice(slice []string, find string) bool {
	for _, v := range slice {
		if v == find {
			return true
		}
	}

	return false
}
```

```go
func TestStrInSlice(t *testing.T) {
	got := StrInSlice([]string{"a", "b"}, "c")

	if got == true {
		t.Errorf("expecting false, got true")
	}
}
```

This function is very simple, and there are only few test cases for it. However, real-world functions need more test cases and table tests are very helpful here:

```go
func TestStrInSlice(t *testing.T) {
	var tests = []struct{
		slice []string
		find  string
		want  bool
	}{
		{[]string{"a", "b"}, "c", false},
		{[]string{"a", "b"}, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.find, func(t *testing.T) {
			got := StrInSlice(tt.slice, tt.find)
			if got != tt.want {
				t.Errorf("expecting %t, got %t", tt.want, got)
			}
		})
	}
}
```

Now let's take more complex code which calls external API and does something with the response. In this example we calculate average stars count per repo of the specified GitHub user:

```go
type Repo struct {
	StargazersCount int `json:"stargazers_count"`
}

func GetAverageStarsPerRepo(username string) (float64, error) {
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
```

And test for it:

```go
func TestGetAverageStarsPerRepo(t *testing.T) {
	var tests = []struct {
		username string
		want     float64
	}{
		{"octocat", 1480.375000},
		{"plutov", 15.566667},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
            got, err := GetAverageStarsPerRepo(tt.username)
            // Don't omit errors even in tests
			if err != nil {
				t.Errorf("expecting nil err, got %v", err)
			}
			if got != tt.want {
				t.Errorf("expecting %f, got %f", tt.want, got)
			}
		})
	}
}
```

It may work well in the beginning, however it's not a good test, it can be flaky, because API may not be available, or testing server has no external connectivity, or simply API response may change (amount of stars).

So how do we call this function in a test, but also avoid testing the HTTP call? We have to restructure our program and make it more modular, create an interface for GitHub API and mock it.

```go
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

	repos := []Repo{}
	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}
```

The `GetAverageStarsPerRepo` function now has to accept the instance of API as the first argument, which can be replaced by Mock in tests:

```go
func GetAverageStarsPerRepo(repositoriesAPI RepositoriesAPI, username string) (float64, error) {
	repos, err := repositoriesAPI.GetRepos(username)
	if err != nil {
		return 0, err
	}

	// ...
}
```

As you can see the function now is much smaller and easier to read. Also tests will be much faster which is very important in bigger complex systems, developers usually don't like to wait long times for their tests to complete (or fail).

And tests would change a bit as well:

```go
// ...

mock := new(Mock)
got, err := GetAverageStarsPerRepo(mock, tt.username)

// ...
```

If we would do this from the beginning, it would save us some time of restructuring the program. That's what I mean when I say "testable code".

Another good practice for testing in Go is to put your tests into a separate `_test` package, this prevents access to private variables, which also allows you to write tests as though you were a real user of the package.

```go
package testable_test

import (
	"testing"

	"github.com/plutov/packagemain/19-testable-go-code/testable"
)

// ...

testable.StrInSlice(...)

// ...
```

There are few more global good practices that can be applied to any language, but we won't go into the details. Such can be:

- Don't use global state, it makes tests difficult to write and make them flaky by default.
- Separate unit tests from integration tests, the latter one doesn't use Mocks and is slower.

And yes, testable code is definitely a good code!