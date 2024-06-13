# Performance Benchmarking: gRPC+Protobuf vs. HTTP+JSON

## Generate unimplemented server gRPC client stub

```
protoc -I./grpc  --go_out=. --go-grpc_out=. \
--go-grpc_opt=require_unimplemented_servers=true  \
--validate_out="lang=go:." users.proto
```

## Run benchmarks

```
go test -bench=. -benchmem=1  -benchtime=30s
```