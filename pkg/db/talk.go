package db

import (
	"github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type MessageData struct {
	Id          string    `db:"id"`
	To          string    `db:"to_"`
	From        string    `db:"from_"`
	ContentType int       `db:"contentType"`
	Text        string    `db:"text"`
	Metadata    string    `db:"metadata"`
	CreatedAt   time.Time `db:"createdAt"`
	UpdatedAt   time.Time `db:"updatedAt"`
}

type InsertMessageReq struct {
	Id          string
	To          string
	From        string
	ContentType protobuf.MessageType
	Text        string
	Metadata    string
}

func GetMessage(messageId string) (*protobuf.Message, error) {
	PingAndReconnect()

	var messageData MessageData
	row := MsngerDB.QueryRow(`select id, to_, from_, contentType, text, metadata, createdAt, updatedAt from msnger.Message where id = ? limit 1`, messageId)
	err := row.Scan(&messageData.Id, &messageData.To, &messageData.From, &messageData.ContentType, &messageData.Text, &messageData.Metadata, &messageData.CreatedAt, &messageData.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// todo : metadata

	msg := protobuf.Message{
		Id:          messageData.Id,
		To:          messageData.To,
		From:        messageData.From,
		ContentType: protobuf.MessageType(messageData.ContentType),
		Text:        messageData.Text,
		Metadata:    nil,
		CreatedAt:   timestamppb.New(messageData.CreatedAt),
		UpdatedAt:   timestamppb.New(messageData.UpdatedAt),
	}

	return &msg, nil
}

func InsertMessage(req *InsertMessageReq) error {
	ins, err := MsngerDB.Prepare(`INSERT INTO msnger.Message (id, to_, from_, contentType, text, metadata) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = ins.Exec(req.Id, req.To, req.From, req.ContentType, req.Text, req.Metadata)
	return err
}
