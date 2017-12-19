package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/plutov/go-youtube-channel/01-grpc/proto"
	"github.com/plutov/go-youtube-channel/01-grpc/server/blockchain"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to listen port 8080: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterBlockchainServer(srv, &Server{
		Blockchain: blockchain.NewBlockchain(),
	})
	srv.Serve(listener)
}

// Server implements proto.BlockchainServer interface
type Server struct {
	Blockchain *blockchain.Blockchain
}

// AddBlock : adds new block to blockchain
func (s *Server) AddBlock(ctx context.Context, in *proto.AddBlockRequest) (*proto.AddBlockResponse, error) {
	block := s.Blockchain.AddBlock(in.Data)
	return &proto.AddBlockResponse{
		Hash: string(block.Hash),
	}, nil
}

// GetBlockchain : returns blockchain
func (s *Server) GetBlockchain(ctx context.Context, in *proto.GetBlockchainRequest) (*proto.GetBlockchainResponse, error) {
	resp := new(proto.GetBlockchainResponse)
	for _, b := range s.Blockchain.Blocks {
		resp.Blocks = append(resp.Blocks, &proto.Block{string(b.PrevBlockHash), string(b.Data), string(b.Hash)})
	}

	return resp, nil
}
