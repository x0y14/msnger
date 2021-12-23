package main

import (
	"context"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/op"
	"io"
	"log"
)

func main() {
	cl := *op.CreateClient("localhost:9191")

	ops, err := cl.FetchOps(context.Background(), &pb.FetchOpsRequest{LastRevisionId: 1})
	if err != nil {
		log.Fatalf("failed to polling ops: %v", err)
	}
	for {
		operation, err := ops.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive op from ops: %v", err)
		}
		log.Printf("%v", operation.String())
	}
}
