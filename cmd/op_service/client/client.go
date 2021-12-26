package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/account"
	"github.com/x0y14/msnger/pkg/service/op"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
)

func main() {
	var email string
	var password string
	var userId string
	var token string

	flag.StringVar(&email, "e", "", "email for login")
	flag.StringVar(&password, "p", "", "password for login")
	flag.Parse()

	// login
	acCl := *account.CreateClient("localhost:9191")
	loginResult, err := acCl.Login(context.Background(), &protobuf.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}
	userId = loginResult.UserId
	token = loginResult.Jwt

	log.Printf("Msnger >> LoginSuccess\n")
	fmt.Printf("Your UserId: %v\n", userId)
	fmt.Printf("Your Token: %v\n", token)

	// fetch ops
	var lastRevisionId uint64 = 0
	opCl := *op.CreateClient("localhost:9292")

	bearer := "Bearer " + token
	md := metadata.Pairs("authorization", bearer)
	ctxAuthed := metadata.NewOutgoingContext(context.Background(), md)

	//ctxAuthed := context.Background()

	stream, err := opCl.FetchOps(ctxAuthed, &protobuf.FetchOpsRequest{LastRevisionId: lastRevisionId})
	if err != nil {
		log.Fatalf("failed to fetch op: %v", err)
	}

	for {
		operation, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive op from stream: %v", err)
		}
		//if operation.Type == protobuf.OperationType_SEND_MESSAGE_RECV {
		//	log.Printf("Receive Message: %v\n", operation.Message.Text)
		//}
		//log.Printf("GOT OP: %v\n", operation.String())
		op.ShowOp(userId, operation)
		lastRevisionId = operation.RevisionId
	}
}
