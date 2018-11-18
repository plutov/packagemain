## Code review of https://github.com/contraband/anderson

1. Fork
2. `git clone git@github.com:plutov/anderson.git`, I clone to original org name of the repo, so Go imports will still work.

### Package structure recommendations

`main.go` imports local package `anderson/anderson`, and for me it seems like a duplication in naming. It is a good practice to have `main.go` files in `cmd` folder with subfolders for each executable. And also we can move files from anderson into root folder.

1. Move main.go file, don't forget to change README and `script/ci`, `go get -u github.com/contraband/anderson/cmd/anderson`.
2. Move files from `./anderson` to `.`
3. Use go modules, `go init` will parse Godeps. Then remove Godeps and vendor. Remove `script/deps`, change `script/ci`

```bash
export GO111MODULE=on
go mod init
go mod tidy
```

## Testing


We will test this program on one of my projects.

```bash
go run cmd/anderson/main.go
> Hold still citizen, scanning dependencies for contraband...
> exit status 1%!(EXTRA []interface {}=[])
```

Error is not meaningful.

In `fatalf` add `args...`, also change message `fatalf("could not get dependencies: %s", err.Error())`.

```bash
go run cmd/anderson/main.go
Could not find github.com/cloudfoundry-incubator/candiedyaml in your GOPATH...
```

Right, because I don't have GOPATH, since it's not required now. Go uses default `~/go` as GOPATH, let's add it.

```go
func Gopaths() ([]string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	if gopath == "" {
		return []string{}, errors.New("GOPATH not set")
	}

	return strings.Split(gopath, ":"), nil
}
```