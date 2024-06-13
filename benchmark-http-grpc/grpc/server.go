package grpc_server

import (
	"context"
	"log"
	"net"

	"github.com/plutov/packagemain/benchmark-http-grpc/grpc/gen"
	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedUsersServer
}

func (s *Server) CreateUser(ctx context.Context, user *gen.User) (*gen.CreateUserResponse, error) {
	return &gen.CreateUserResponse{
		Message: "ok",
		Code:    201,
		User:    user,
	}, nil
}

func Start() {
	lis, err := net.Listen("tcp", ":60000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	gen.RegisterUsersServer(srv, &Server{})
	srv.Serve(lis)
}
