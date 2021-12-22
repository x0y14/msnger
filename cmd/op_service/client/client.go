package main

import (
	"context"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func createClient(addr string) *pb.OpServiceClient {

	creds := insecure.NewCredentials()

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect server: %v : %v", addr, err)
	}

	client := pb.NewOpServiceClient(conn)
	return &client
}

func main() {
	cl := *createClient("localhost:9191")

	ops, err := cl.FetchOps(context.Background(), &pb.FetchOpsRequest{LastRevisionId: 1})
	if err != nil {
		log.Fatalf("failed to polling ops: %v", err)
	}
	for {
		op, err := ops.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive op from ops: %v", err)
		}
		log.Printf("%v", op.String())
	}
}
