package main_test

import (
	"context"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/account"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"
)

const (
	testUserEmail    = "tom@example.com"
	testUserPassword = "12345"
	testUserId       = "Uc7401iqs1s4362f70sa0"
	testUserToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJVYzc0MDFpcXMxczQzNjJmNzBzYTAifQ.5se5p4_FghHJod_JnzNo5AzTqHWiVJpSm1myMxQNzSA"
)

func TestCreateAccount(t *testing.T) {
	cl := *account.CreateClient("localhost:9191")

	result, err := cl.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
	})
	if err != nil {
		t.Fatalf("failed to create an account: %v", err)
	}

	log.Printf("UserId: %v\n", result.GetUserId())
	log.Printf("Token: %v\n", result.GetJwt())
}

func TestLogin(t *testing.T) {
	cl := *account.CreateClient("localhost:9191")
	result, err := cl.Login(context.Background(), &pb.LoginRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
	})
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	log.Printf("UserId: %v\n", result.GetUserId())
	log.Printf("Token: %v\n", result.GetJwt())
}

func TestGetEmail_NoLogin(t *testing.T) {
	cl := *account.CreateClient("localhost:9292")
	result, err := cl.GetEmail(context.Background(), &pb.GetEmailRequest{})
	if err != nil {
		t.Fatalf("failed to get email: %v", err)
	}
	log.Printf("Email: %v", result.Email)
}

func TestGetEmail_Login(t *testing.T) {
	cl := *account.CreateClient("localhost:9292")

	bearer := "Bearer " + testUserToken

	md := metadata.Pairs("authorization", bearer)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	result, err := cl.GetEmail(ctx, &pb.GetEmailRequest{})
	if err != nil {
		t.Fatalf("failed to get email: %v", err)
	}
	log.Printf("Email: %v", result.Email)
}
