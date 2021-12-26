package db_test

import (
	"fmt"
	"github.com/x0y14/msnger/pkg/db"
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
	"testing"
)

func TestInsertMessage(t *testing.T) {
	msgId := misc.GenerateMessageId()
	req := db.InsertMessageReq{
		Id:          msgId,
		To:          msgId,
		From:        msgId,
		ContentType: protobuf.MessageType_TEXT,
		Text:        "hello",
		Metadata:    "{}",
	}
	err := db.InsertMessage(&req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMessage(t *testing.T) {
	msgId := "Mc73h5mqs1s419frhpoa0"
	msg, err := db.GetMessage(msgId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Id: %v\n", msg.Id)
	fmt.Printf("To: %v\n", msg.To)
	fmt.Printf("From: %v\n", msg.From)
	fmt.Printf("ContentType: %v\n", msg.ContentType)
	fmt.Printf("Text: %v\n", msg.Text)
	fmt.Printf("Metadata: %v\n", msg.Metadata)
	fmt.Printf("CreatedAt: %v\n", msg.CreatedAt.String())
	fmt.Printf("UpdatedAt: %v\n", msg.UpdatedAt.String())
}
