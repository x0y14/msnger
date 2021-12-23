package main

import (
	"context"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/account"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	cl := *account.CreateClient("localhost:9292")

	// login必要ない
	resultOfLogin, err := cl.Login(context.Background(), &pb.LoginRequest{Email: "sample@example.com", Password: "p@ssword"})
	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}
	myUserId := resultOfLogin.GetUserId()
	myToken := resultOfLogin.GetJwt()
	log.Printf("I got UserId: %v, Token: %v\n", myUserId, myToken)

	// 失敗: loginなし
	resultOfGetEmail, err := cl.GetEmail(context.Background(), &pb.GetEmailRequest{})
	if err != nil {
		log.Printf("failed to get email: %v", err)
	}

	// 成功: loginあり
	bearer := "Bearer " + myToken
	md := metadata.New(map[string]string{"authorization": bearer})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resultOfGetEmail, err = cl.GetEmail(ctx, &pb.GetEmailRequest{})
	if err != nil {
		log.Fatalf("failed to get email: %v", err)
	}
	log.Printf("success to get email: %v\n", resultOfGetEmail.GetEmail())
}
