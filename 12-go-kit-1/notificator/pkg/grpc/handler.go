package grpc

import (
	"context"
	"errors"

	grpc "github.com/go-kit/kit/transport/grpc"
	endpoint "github.com/plutov/packagemain/12-go-kit-1/notificator/pkg/endpoint"
	pb "github.com/plutov/packagemain/12-go-kit-1/notificator/pkg/grpc/pb"
	context1 "golang.org/x/net/context"
)

// makeSendEmailHandler creates the handler logic
func makeSendEmailHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.SendEmailEndpoint, decodeSendEmailRequest, encodeSendEmailResponse, options...)
}

// decodeSendEmailResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain sum request.
// TODO implement the decoder
func decodeSendEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Notificator' Decoder is not impelemented")
}

// encodeSendEmailResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeSendEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Notificator' Encoder is not impelemented")
}
func (g *grpcServer) SendEmail(ctx context1.Context, req *pb.SendEmailRequest) (*pb.SendEmailReply, error) {
	_, rep, err := g.sendEmail.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SendEmailReply), nil
}
