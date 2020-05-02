## packagemain #19: Writing testable Go code

When I say "testable code", what I mean is code that can be easily programmatically verified. We can say that code is testable when we don't have to change the code itself when we're adding a unit test to it. It doesn't matter if you're following test-driven development or not, testable code makes your program more flexible and maintainable, due to its modularity.

Go has robust built-in testing functionality, so in most cases you don't need to import any third-party testing packages. So start clean in the beginning, and if it's not enough, you can add helper package ([assert](https://pkg.go.dev/github.com/stretchr/testify/assert) for example) later.

### SOLID

First of all, understanding of SOLID principles will help you with writing testable code. I won't go into details, but Single Responsibility and Dependency Inversion will help you a lot.

For example, it's much easier and cleaner to test small function which does only one thing. For example, function `strInSlice` is perfectly testable function:

```go
func strInSlice(slice []string, find string) bool {
	for _, v := range slice {
		if v == find {
			return true
		}
	}

	return false
}
```

```go
func Test_strInSlice(t *testing.T) {
	got := strInSlice([]string{"a", "b"}, "c")

	if got == true {
		t.Errorf("expecting false, got true")
	}
}
```

This function is very simple, and there are only few test cases for it. However, real-world functions need more test cases and table test are very helpful here:

```go
func Test_strInSlice(t *testing.T) {
	var tests = []struct {
		slice []string
		find  string
		want  bool
	}{
		{[]string{"a", "b"}, "c", false},
		{[]string{"a", "b"}, "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.find, func(t *testing.T) {
			got := strInSlice(tt.slice, tt.find)
			if got != tt.want {
				t.Errorf("expecting %t, got %t", tt.want, got)
			}
		})
	}
}
```

Now let's take more complex code which calls external API and checks the response:

```go
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
```

And test for it:

```go
func Test_getAveragePerRepo(t *testing.T) {
	var tests = []struct {
		username string
		want     float64
	}{
		{"octocat", 1480.375000},
		{"plutov", 15.566667},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			got, err := getAveragePerRepo(tt.username)
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

2. mocks
3. net/http/httptest
5. separate pkg_test package
