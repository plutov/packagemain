package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/net/context"

	"github.com/plutov/packagemain/01-grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	client := proto.NewBlockchainClient(conn)

	addResp, addErr := client.AddBlock(context.Background(), &proto.AddBlockRequest{
		Data: "test",
	})
	if addErr != nil {
		log.Fatalf("unable to add block: %v", addErr)
	}
	spew.Dump(addResp)

	getResp, getErr := client.GetBlockchain(context.Background(), &proto.GetBlockchainRequest{})
	if getErr != nil {
		log.Fatalf("unable to get blockchain: %v", getErr)
	}
	spew.Dump(getResp)
}
