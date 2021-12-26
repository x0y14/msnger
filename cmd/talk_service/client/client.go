package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/account"
	"github.com/x0y14/msnger/pkg/service/talk"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
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

	fmt.Printf("set message receiver >> ")
	var receiverUserId string
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		if err := stdin.Err(); err != nil {
			log.Fatalf("%v", err)
		}
		receiverUserId = stdin.Text()
	}

	talkCl := *talk.CreateClient("localhost:9090")
	bearer := "Bearer " + token
	md := metadata.Pairs("authorization", bearer)
	ctxAuthed := metadata.NewOutgoingContext(context.Background(), md)

	fmt.Printf("--- [ user input (to %v) ] ---\n", receiverUserId)
	var text string
	for {
		fmt.Printf("you >> ")
		if stdin.Scan() {
			text = stdin.Text()
			_, err = talkCl.SendMessage(ctxAuthed, &protobuf.SendMessageRequest{Message: &protobuf.Message{
				Id:          "",
				From:        userId,
				To:          receiverUserId,
				ContentType: protobuf.MessageType_TEXT,
				Text:        text,
				Metadata:    nil,
				CreatedAt:   nil,
				UpdatedAt:   nil,
			}})
			fmt.Printf("\t SENT = %v\n", text)
		} else {
			log.Fatalf("failed to read line: %v", stdin.Err())
		}
	}
}
