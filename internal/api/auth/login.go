package auth

import (
	"context"

	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login Логин.
func (c Controller) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.FailedPrecondition, "request is nil")
	}
	if req.Username == "" {
		return nil, status.Error(codes.FailedPrecondition, "username is empty")
	}
	if req.Password == "" {
		return nil, status.Error(codes.FailedPrecondition, "password is empty")
	}

	token, err := c.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{RefreshToken: token}, nil
}
