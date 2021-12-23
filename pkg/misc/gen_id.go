package misc

import (
	"github.com/godruoyi/go-snowflake"
	"github.com/rs/xid"
)

func GenerateRevisionId() uint64 {
	id := snowflake.ID()
	return id
}

func GenerateUserId() string {
	return "U" + xid.New().String()
}

func GenerateMessageId() string {
	return "M" + xid.New().String()
}

func GenerateGroupId() string {
	return "G" + xid.New().String()
}
