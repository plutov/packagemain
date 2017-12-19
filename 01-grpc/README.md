#### Goal

Create gRPC Server and Client in Go. We will create a fictional blockchain service.

#### Steps

 - Install [protoc](https://github.com/google/protobuf/releases) compiler.
 - Install `protoc-gen-go plugin`: `go get -u github.com/golang/protobuf/protoc-gen-go`
 - Define service definition in `.proto` file.
 - Build Go gRPC bindings from `.proto` file. `protoc --go_out=plugins=grpc:. proto/service.proto`
 - Install grpc Go package - `go get -u google.golang.org/grpc`.
 - Install context package - `go get -u golang.org/x/net/context`.
 - Install protobuf package - `go get -u github.com/golang/protobuf/proto`
 - Implement Server, interface `BlockchainServer`.
 - Create a client using `BlockchainClient`.
 - Run server first.
 - Run client.

#### Comments

 - Use [go-spew](https://github.com/davecgh/go-spew) to dump structs. 

#### Resources

 - https://jeiwan.cc/posts/building-blockchain-in-go-part-1/