package account

import (
	"context"
	"github.com/x0y14/msnger/pkg/auth"
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceServer struct {
	protobuf.UnimplementedAccountServiceServer
}

func (s *ServiceServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if fullMethodName == "/account.AccountService/CreateAccount" {
		return ctx, nil
	} else if fullMethodName == "/account.AccountService/Login" {
		return ctx, nil
	}

	ctx, err := auth.Authentication(ctx)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (s *ServiceServer) CreateAccount(ctx context.Context, req *protobuf.CreateAccountRequest) (*protobuf.CreateAccountResult, error) {
	newUserId := misc.GenerateUserId()
	token, err := auth.GenerateJWTToken(newUserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token.")
	}
	return &protobuf.CreateAccountResult{
		UserId: newUserId,
		Jwt:    token,
	}, nil
}

func (s *ServiceServer) Login(ctx context.Context, req *protobuf.LoginRequest) (*protobuf.LoginResult, error) {
	if req.Email != "sample@example.com" || req.Password != "p@ssword" {
		return nil, status.Errorf(codes.Unauthenticated, "failed to login, email or password was invalid")
	}
	sampleUserId := misc.GenerateUserId()
	token, err := auth.GenerateJWTToken(sampleUserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token.")
	}
	return &protobuf.LoginResult{
		UserId: sampleUserId,
		Jwt:    token,
	}, nil
}

func (s *ServiceServer) GetEmail(ctx context.Context, req *protobuf.GetEmailRequest) (*protobuf.GetEmailResult, error) {
	return &protobuf.GetEmailResult{Email: "sample@example.com"}, nil
}
