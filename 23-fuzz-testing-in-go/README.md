## packagemain #23: Fuzz Testing in Go

### What is fuzzing

Fuzzing or fuzz testing is a method of giving random unexpected input to your programs to test for possible crashes or edge cases. Fuzzing can shed a light on some logical bugs or performance problems, so it's always worth adding to a code where stability and performance matter.

### Go projects for fuzzing

Currently, there are few well-supported projects to do fuzzing in go:
- [go-fuzz](https://github.com/dvyukov/go-fuzz)
- and [gofuzz](https://github.com/google/gofuzz)

But we're not going to review them in this video since we have some great news, as Go team has accepted a proposal to add fuzz testing support to the language. It will be available in the standard toolchain in Go 1.18 - [docs](https://go.dev/doc/fuzz/).

### Install Go 1.18

At the moment of writing this post Go 1.18 is only in beta, so let's install it first with a help of `gotip`.

```
go install golang.org/dl/gotip@latest
gotip download
```

Note: probably when you read it Go 1.18 is already released, so you can upgrade your go command in a regular way.

### Our "pet" function

Let's write a simple function for which we'll add fuzzing later on. And we will add some errors intentionally.

The function `Equal` will compare two slices of bytes element by element.

```go
package fuzztestingingo

func Equal(a []byte, b []byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
```

### Writing fuzz test

1. Create a file `equal_test.go`
2. Let's include a simple regular test

```go
package fuzztestingingo

import "testing"

func TestEqual(t *testing.T) {
	if !Equal([]byte{'f', 'u', 'z', 'z'}, []byte{'f', 'u', 'z', 'z'}) {
		t.Error("expected true, got false")
	}
}
```

```
gotip test .
ok github.com/plutov/packagemain/23-fuzz-testing-in-go	0.922s
```

The test works but it checks only a simple use case, and as you probably already noticed our function has some edge cases. Let's try to uncover them by writing a fuzz test.

Fuzz tests can be included in your regular `_test.go` files using the functions that start with `Fuzz` that accept new type `*testing.F`.

```go
func FuzzEqual(f *testing.F) {
    // target, can be only one per test
	// values of a and b will be auto-generated
	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		Equal(a, b)
	})
}
```

Note: There can be only one target per test.

To enable fuzzing, we have to run `go test` with the `-fuzz` flag:

```
gotip test -fuzz .
fuzz: elapsed: 0s, gathering baseline coverage: 0/2 completed
failure while testing seed corpus entry: FuzzEqual/84ed65595ad05a58e293dbf423c1a816b697e2763a29d7c37aa476d6eef6fd60
fuzz: elapsed: 0s, gathering baseline coverage: 1/2 completed
--- FAIL: FuzzEqual (0.02s)
    --- FAIL: FuzzEqual (0.00s)
        testing.go:1349: panic: runtime error: index out of range [0] with length 0
```

### Fixing the "pet" function

We found our error, as our code doesn't check the size of the slice. Let's fix that and run the fuzz again.

```go
package fuzztestingingo

func Equal(a []byte, b []byte) bool {
    // if length is not the same - slices are not equal
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
```

```
gotip test -fuzz .
fuzz: elapsed: 0s, gathering baseline coverage: 0/11 completed
fuzz: elapsed: 0s, gathering baseline coverage: 11/11 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 542957 (180982/sec), new interesting: 0 (total: 11)
fuzz: elapsed: 6s, execs: 1035678 (164216/sec), new interesting: 0 (total: 11)
...
```

It is up to you to decide how long to run fuzzing. It is very possible that execution of fuzzing could run indefinitely if it doesn’t find any errors, like in our case. We can add `-fuzztime` argument, that tells how many iterations to run.

```
gotip test -fuzz=. -fuzztime=5s .
fuzz: elapsed: 0s, gathering baseline coverage: 0/10 completed
fuzz: elapsed: 0s, gathering baseline coverage: 10/10 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 474778 (158251/sec), new interesting: 0 (total: 10)
fuzz: elapsed: 5s, execs: 729255 (121223/sec), new interesting: 0 (total: 10)
PASS
ok  	github.com/plutov/packagemain/23-fuzz-testing-in-go	5.557s
```

Now let's review the output of the fuzzing, there are multiple metrics:

- elapsed: the amount of time that has elapsed since the process began
- execs: the total number of inputs that have been run against the fuzz target
- new interesting: the total number of “interesting” inputs that have been added to the generated corpus during this fuzzing execution

Please be aware that fuzzing can consume a lot of memory and may impact your machine’s performance while it runs, so you should be careful running fuzz in your CI pipeline.
