package auth_test

import (
	"context"
	"github.com/x0y14/msnger/pkg/auth"
	"log"
	"testing"
)

func TestAuthentication(t *testing.T) {
	userId := "aaaa4"
	jwt, err := auth.GenerateJWTToken(userId)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("jwt: %v\n", jwt)

	bearer := "Bearer " + jwt
	ctx := context.WithValue(context.Background(), "authorization", bearer)

	ctxIncludedUserId, err := auth.Authentication(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("expected: %v, actual: %v, isSame: %v\n", userId, ctxIncludedUserId.Value("userId"), userId == ctxIncludedUserId.Value("userId"))
}
