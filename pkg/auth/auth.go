package auth

import (
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ServerSecret = "Kishibe Rohan said 'i refuse.'"
)

type UserInfo struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func GenerateJWTToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserInfo{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{},
	})
	tokenStr, err := token.SignedString([]byte(ServerSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func Authentication(ctx context.Context) (context.Context, error) {
	tokenStr, err := grpcAuth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, fmt.Errorf("failed to pick up token: %v", err)
	}

	userInfo := UserInfo{}
	token, err := jwt.ParseWithClaims(tokenStr, &userInfo, func(token *jwt.Token) (interface{}, error) {
		return []byte(ServerSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "jwt was invalid.")
	}

	ctxIncludedUserId := context.WithValue(ctx, "userId", userInfo.UserId)
	return ctxIncludedUserId, nil
}
