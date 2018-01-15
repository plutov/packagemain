#### Goal

Create gRPC Server and Client in Go. We will create a fictional blockchain service.

#### Steps

 - Install [protoc](https://github.com/google/protobuf/releases) compiler.
 - Install `protoc-gen-go plugin`: `go get -u github.com/golang/protobuf/protoc-gen-go`
 - Define service definition in `.proto` file.
 - Build Go bindings from `.proto` file. `protoc --go_out=plugins=grpc:. proto/blockchain.proto`
 - Install grpc Go package - `go get -u google.golang.org/grpc`.
 - Install context package - `go get -u golang.org/x/net/context`.
 - Install protobuf package - `go get -u github.com/golang/protobuf/proto`
 - Implement Server, interface `BlockchainServer`.
 - Create a client using `BlockchainClient`.
 - Run server first.
 - Run client.

#### Usage

Start server:
```
go run server/main.go
```

Add block as client:
```
go run client/main.go --add
```

get blockchain as client:
```
go run client/main.go --list
```
