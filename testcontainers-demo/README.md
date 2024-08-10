## Emulating real dependencies in Integration Tests using Testcontainers\

Run integration tests using Testcontainers:

```bash
go test -tags=integration -v ./...
```

Run unit tests:

```bash
go test -v ./...
```

Run integration tests with real dependencies:

```bash
docker-compose up -d
go test -tags=integration -v ./...
```
