package auth

import (
	"context"
	"errors"
	auth1 "github.com/tgkzz/auth/gen/go/auth"
	"github.com/tgkzz/auth/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	auth1.UnimplementedAuthServiceServer
	auth Auth
}

type Auth interface {
	Register(
		ctx context.Context,
		username string,
		password string,
	) (userId int64, err error)

	Login(
		ctx context.Context,
		username string,
		password string,
	) (token string, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	auth1.RegisterAuthServiceServer(gRPCServer, &serverApi{auth: auth})
}

func (s *serverApi) Register(ctx context.Context, req *auth1.RegisterRequest) (*auth1.RegisterResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username must not be empty")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password must not be empty")
	}

	uid, err := s.auth.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &auth1.RegisterResponse{UserId: uid}, nil
}

func (s *serverApi) Login(ctx context.Context, req *auth1.LoginRequest) (*auth1.LoginResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username must not be empty")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password must not be empty")
	}

	jwtToken, err := s.auth.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &auth1.LoginResponse{Token: jwtToken}, nil
}
