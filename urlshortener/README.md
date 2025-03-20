### unit tests

```bash
mockery --all --with-expecter
go test -v ./...
```

### integration tests

```bash
go test -tags=integration -v ./...
```

