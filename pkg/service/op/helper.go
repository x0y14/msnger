package op

import (
	"github.com/x0y14/msnger/pkg/protobuf"
	"log"
)

func ShowOp(receiverId string, operation *protobuf.Operation) {
	switch operation.Type {
	case protobuf.OperationType_SEND_MESSAGE_SEND:
		log.Printf("SendMessage (own): {ID: %v, FROM: %v, TO: %v, TEXT: %v}\n", operation.Message.Id, operation.Message.From, operation.Message.To, operation.Message.Text)
	case protobuf.OperationType_SEND_MESSAGE_RECV:
		if operation.Message.From == receiverId {
			log.Printf("ReceiveMessage (own): {ID: %v, FROM: %v, TO: %v, TEXT: %v}\n", operation.Message.Id, operation.Message.From, operation.Message.To, operation.Message.Text)
		} else {
			log.Printf("ReceiveMessage: {ID: %v, FROM: %v, TO: %v, TEXT: %v}\n", operation.Message.Id, operation.Message.From, operation.Message.To, operation.Message.Text)
		}
	default:
		log.Printf("Receive Op: %v\n", operation.String())
	}
}
