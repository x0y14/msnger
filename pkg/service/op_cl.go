package service

import (
	"context"
	"github.com/x0y14/msnger/pkg/protobuf"
	"github.com/x0y14/msnger/pkg/service/op"
	"google.golang.org/grpc/metadata"
)

var OpCl protobuf.OpServiceClient

func init() {
	OpCl = *op.CreateClient("localhost:9292")
}

func CreateAdminCtx() context.Context {
	bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJVYzczdTA1cXMxczQybWJiOTVybjAifQ.kLU2fNFuA_tXf36EAFuJq8NxX3F7Z49uPWZ0PE0OEMA"
	md := metadata.Pairs(
		"authorization", bearer,
		"userId", "Uc73u05qs1s42mbb95rn0")
	return metadata.NewOutgoingContext(context.Background(), md)
}
