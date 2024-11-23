package auth

import (
	"context"
	"errors"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login Логин.
func (c Controller) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is empty")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	token, err := c.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrorInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{RefreshToken: token}, nil
}
