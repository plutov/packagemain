## packagemain #18: Basic CI pipeline for Go project using GitHub Actions

Today we'll talk about [GitHub Actions](https://github.com/features/actions), which gives free CI/CD tools for Open Source projects.

GitHub Actions are based on the concept of Workflows, that define the execution order and flow.

In this video we'll build a dummy Go program with a very basic pipeline that after each Pull Request or push to master will run lint and unit tests.

And when new tag is created in the repository, it will create a new GitHub release with artifacts attached.

### Go program

For the sake of simplicity we'll create a Fibonacci counter.

#### pkg/fib/fib.go

```go
package fib

// Fibonacci ..
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
```

And some table tests for this function:

#### pkg/fib/fib_test.go

```go
package fib

import "testing"

func TestFibonacci(t *testing.T) {
	data := []struct {
		n    int
		want int
	}{
		{0, 0}, {1, 1}, {2, 1}, {3, 2}, {4, 3}, {5, 5}, {6, 8}, {10, 55},
	}

	for _, d := range data {
		if got := Fibonacci(d.n); got != d.want {
			t.Errorf("Invalid Fibonacci value for N: %d, got: %d, want: %d", d.n, got, d.want)
		}
	}
}
```

Let's make this code executable by calling `Fibonacci` function from main package.

#### cmd/main.go

```go
package main

import (
	"flag"
	"fmt"

	"github.com/plutov/github-actions-go/pkg/fib"
)

func main() {
	n := flag.Int("n", 1, "fibonacci input number")
	flag.Parse()

	fmt.Printf("result: %d", fib.Fibonacci(*n))
}
```

### Test Pipeline

As said previously, GitHub Actions are based on the concept of Workflows. A workflow is a set of jobs and steps that are executed when some condition or event is met, for example tag or PR is created.

It is possible to have multiple workflows by simply creating separate files in `.github/workflows` folder.

One more amazing thing is that you don't need to reinvent the wheel and you can use actions built by GitHub or community.

Here is how our Test workflow looks:

#### .github/workflows/test.yml

```yaml
on: [pull_request]
name: Test
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.11, 1.12, 1.13]
        platform: [ubuntu-latest]
    runs-on: ${{ matrix.platform }}
    steps:
    - name: Install Go
      uses: actions/setup-go@v1
      with:
        go-version: ${{ matrix.go-version }}

    - name: Checkout code
      uses: actions/checkout@v1

    - name: Run Tests
      run: go test -race -v ./pkg/...
```

I love matrix build feature, which is a set of different configurations of the runner environment. In our case we will run tests on few Go versions to ensure compatibility.

### Release Pipeline

#### .github/workflows/release.yml

```yaml
on:
  create:
    tags:
      - v*
name: Release
jobs:
  release:
    name: Release on GitHub
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v1

    - name: Validates GO releaser config
      uses: docker://goreleaser/goreleaser:latest
      with:
        args: check

    - name: Create release on GitHub
      uses: docker://goreleaser/goreleaser:latest
      with:
        args: release
      env:
        GITHUB_TOKEN: ${{secrets.GITHUB_TOKEN}}
```

Here we're using [GoReleaser](https://goreleaser.com/) official Docker image, which will publish release with artifacts directly to GitHub using `GITHUB_TOKEN`. You can manage `secrets` on your repository's settings page, but in case if `GITHUB_TOKEN` you don't need to do anything as this variable is injected automatically by the Actions platform.

```
goreleaser init
```

### Push and Release

Now as our code and workflows are ready, we can try to push it to our remote repository and check if it works.

```
git co -b workflows
git add -A
git commit -m "GitHub Workflows"
git push origin workflows
```

After PR is merged, let's create a tag:

```
git tag v1.0.0
git push --tags
```

### Conclusion

That was it for setting up a basic workflows for Go project, but you can do much more with Actions. Check out [Actions Marketplace](https://github.com/marketplace?type=actions), where you can find Actions of your interest, I attached few helpful links below.