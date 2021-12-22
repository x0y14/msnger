package main

import (
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/op"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	addr := "localhost:9191"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOpServiceServer(s, &op.ServiceServer{})
	log.Printf("op service server listening at: %v", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
