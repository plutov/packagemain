## Emulating real dependencies in Integration Tests using Testcontainers

We built a dead simple URL shortener API in Go with 2 endpoints:
- `/create?url=` - shortens the given URL.
- `/get?key=` - returns the original URL for a given key.

It has 2 dependencies:
- MongoDB as the database.
- Redis for caching.

### Run unit tests

```bash
go test -v ./...
```

### Generate mocks for unit tests

Make sure to install [mockery](https://github.com/vektra/mockery) first.

```bash
mockery --all --with-expecter
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
