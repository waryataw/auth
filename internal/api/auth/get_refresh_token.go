package auth

import (
	"context"

	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken Получить Refresh Токен.
func (c Controller) GetRefreshToken(ctx context.Context, req *authv1.GetRefreshTokenRequest) (*authv1.GetRefreshTokenResponse, error) {
	if req == nil || req.OldRefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "old refresh token is required")
	}

	token, err := c.service.UpdateRefreshToken(ctx, req.OldRefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.GetRefreshTokenResponse{RefreshToken: token}, nil
}
