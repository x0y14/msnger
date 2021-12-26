package talk

import (
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func CreateClient(addr string) *pb.TalkServiceClient {

	creds := insecure.NewCredentials()

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect server: %v : %v", addr, err)
	}

	client := pb.NewTalkServiceClient(conn)
	return &client
}
