package auth

import (
	"context"

	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccessToken Получить Access Токен.
func (c Controller) GetAccessToken(ctx context.Context, req *authv1.GetAccessTokenRequest) (*authv1.GetAccessTokenResponse, error) {
	if req == nil || req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	token, err := c.service.GetAccessToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.GetAccessTokenResponse{AccessToken: token}, nil
}
