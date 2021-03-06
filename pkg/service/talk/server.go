package talk

import (
	"context"
	"github.com/x0y14/msnger/pkg/db"
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type ServiceServer struct {
	protobuf.UnimplementedTalkServiceServer
}

func (s *ServiceServer) SendMessage(ctx context.Context, req *protobuf.SendMessageRequest) (*protobuf.Message, error) {
	msgId := misc.GenerateMessageId()
	sendMessageOp := &protobuf.SendOpRequest{Op: &protobuf.Operation{
		RevisionId: 0,
		Type:       protobuf.OperationType_SEND_MESSAGE_SEND,
		Param1:     req.Message.From,
		Param2:     req.Message.To,
		Param3:     msgId,
		Message:    req.Message,
		CreatedAt:  nil,
		UpdatedAt:  nil,
	}}

	err := db.InsertMessage(&db.InsertMessageReq{
		Id:          msgId,
		To:          req.Message.To,
		From:        req.Message.From,
		ContentType: req.Message.ContentType,
		Text:        req.Message.Text,
		Metadata:    "{}",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert message to db.")
	}

	_, err = service.OpCl.SendOp(service.CreateAdminCtx(), sendMessageOp)
	if err != nil {
		log.Printf("[Talk] SendMessage Err: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send op by talk service.")
	}
	return &protobuf.Message{}, nil
}

func (s *ServiceServer) SendReadReceipt(ctx context.Context, req *protobuf.SendReadReceiptRequest) (*protobuf.SendReadReceiptResult, error) {
	return &protobuf.SendReadReceiptResult{}, nil
}
