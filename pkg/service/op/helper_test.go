package op_test

import (
	"github.com/x0y14/msnger/pkg/service/op"
	"log"
	"testing"
)

func TestGenerateRevisionId(t *testing.T) {
	revisionId := op.GenerateRevisionId()
	log.Printf("generated revisionId (uint64, snowflake): %v", revisionId)
}

func TestGenerateUserId(t *testing.T) {
	userId := op.GenerateUserId()
	log.Printf("generated userId (string, xid): %v", userId)
}

func TestGenerateMessageId(t *testing.T) {
	msgId := op.GenerateMessageId()
	log.Printf("generated msgId (string, xid): %v", msgId)
}
