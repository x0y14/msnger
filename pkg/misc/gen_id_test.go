package misc

import (
	"log"
	"testing"
)

func TestGenerateRevisionId(t *testing.T) {
	revisionId := GenerateRevisionId()
	log.Printf("generated revisionId (uint64, snowflake): %v", revisionId)
}

func TestGenerateUserId(t *testing.T) {
	userId := GenerateUserId()
	log.Printf("generated userId (string, xid): %v", userId)
}

func TestGenerateMessageId(t *testing.T) {
	msgId := GenerateMessageId()
	log.Printf("generated msgId (string, xid): %v", msgId)
}
