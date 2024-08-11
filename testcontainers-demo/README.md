## Emulating real dependencies in Integration Tests using Testcontainers

We built a dead simple URL shortener API in Go with 2 endpoints:
- `/create?url=` - shortens the given URL
- `/get?key=` - redirects to the original URL

It has 2 dependencies:
- Postgres for storing the mappings between short and original URLs.
- Redis for caching.

### Run unit tests

```bash
go test -v ./...
```

### Generate mocks for unit tests

Make sure to install [mockery](https://github.com/vektra/mockery) first.

```bash
mockery --all --recursive --case underscore --with-expecter
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
