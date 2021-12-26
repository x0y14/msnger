package main

import (
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/x0y14/msnger/pkg/auth"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/talk"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	addr := "localhost:9090"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcAuth.UnaryServerInterceptor(auth.Authentication)),
		grpc.StreamInterceptor(grpcAuth.StreamServerInterceptor(auth.Authentication)))
	pb.RegisterTalkServiceServer(s, &talk.ServiceServer{})
	log.Printf("talk service server listening at: %v", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
