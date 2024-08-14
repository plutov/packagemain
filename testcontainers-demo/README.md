## Emulating real dependencies in Integration Tests using Testcontainers

We built a dead simple URL shortener API in Go with 2 endpoints:
- `/create?url=` - shortens the given URL.
- `/get?key=` - returns the original URL for a given key.

It has 2 dependencies:
- MongoDB as the database.
- Redis for caching.

### Run unit tests

Uses mocks created with [mockery](https://github.com/vektra/mockery) as dependencies.

See `unit_test.go` file for more implementation.

```bash
mockery --all --with-expecter
go test -v ./...
```

### Run integration tests using Testcontainers

See `integration_test.go` file for more implementation.

The first run may take a while to download the images. But the subsequent runs are almost instant.

```bash
go test -tags=integration -v ./...
```

### Run integration tests with real dependencies

See `realdeps_test.go` file for more implementation. Make sure to start the services before running the tests.

```bash
docker-compose up -d
go test -tags=realdeps -v ./...
```
