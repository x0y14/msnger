package account

import (
	"context"
	"github.com/x0y14/msnger/pkg/auth"
	"github.com/x0y14/msnger/pkg/db"
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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

	// store account data to db
	err := db.InsertAccount(&db.InsertAccountReq{
		Id:       newUserId,
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		log.Printf("[Account] CreateAccount Err: %v\n", err)
		return nil, status.Errorf(codes.Internal, "error")
	}

	// generate token
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
	account, err := db.SelectAccountWithEmail(&db.SelectAccountReq{Email: req.GetEmail()})
	if err != nil {
		log.Printf("[Account] Login Err: %v\n", err)
		return nil, status.Errorf(codes.Internal, "error")
	}

	// アカウント不明
	if account == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found.")
	}

	// 成功, トークン生成
	jwt, err := auth.GenerateJWTToken(account.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token.")
	}

	return &protobuf.LoginResult{
		UserId: account.GetId(),
		Jwt:    jwt,
	}, nil
}

func (s *ServiceServer) GetEmail(ctx context.Context, req *protobuf.GetEmailRequest) (*protobuf.GetEmailResult, error) {
	userId := ctx.Value("userId").(string)
	account, err := db.SelectAccountWithId(&db.SelectAccountReq{Id: userId})
	if err != nil {
		log.Printf("[Account] GetEmail Err: %v\n", err)
		return nil, status.Errorf(codes.Internal, "error")
	}

	if account == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found.")
	}

	return &protobuf.GetEmailResult{Email: account.Email}, nil
}
