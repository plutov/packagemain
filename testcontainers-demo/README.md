## Emulating real dependencies in Integration Tests using Testcontainers

We built a dead simple URL shortener API in Go with 2 endpoints:
- `/create?url=` - shortens the given URL.
- `/get?key=` - returns the original URL for a given key.

It has 2 dependencies:
- MongoDB as the database.
- Redis for caching.

### Run unit tests

Uses mocks created with [mockery](https://github.com/vektra/mockery) as dependencies.

```bash
mockery --all --with-expecter
go test -v ./...
```

### Run integration tests using Testcontainers

```bash
go test -tags=integration -v ./...
```

### Run integration tests with real dependencies

```bash
docker-compose up -d
go test -tags=integration -v ./...
```
