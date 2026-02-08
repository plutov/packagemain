# jemalloc

Watch the [video](https://www.youtube.com/watch?v=SHmJTgjldgg).

Installation:

```
brew install jemalloc
export CGO_CFLAGS="-I/opt/homebrew/include"
export CGO_LDFLAGS="-L/opt/homebrew/lib"
```

Run with Go's garbage collector:

```
go run . -gc
```

Run with jemalloc:
```
go run .
```

Run with jemalloc and disabled garbage collector:
```
GOGC=off go run .
```
