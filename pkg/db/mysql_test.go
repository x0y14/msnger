package db_test

import (
	"github.com/rs/xid"
	"github.com/x0y14/msnger/pkg/db"
	"github.com/x0y14/msnger/pkg/protobuf"
	"testing"
)

func TestStoreOp(t *testing.T) {
	msg := &protobuf.Message{
		Id:          xid.New().String(),
		To:          xid.New().String(),
		From:        xid.New().String(),
		CreatedAt:   nil,
		ContentType: 0,
		Text:        "hello, operation",
		Metadata:    nil,
	}

	op := &protobuf.Operation{
		RevisionId: 0,
		Type:       protobuf.OperationType_NOOP,
		Param1:     "no4",
		Param2:     "p2",
		Param3:     "p3",
		Msg:        msg,
		CreatedAt:  nil,
	}

	db.StoreOp(op)
}
