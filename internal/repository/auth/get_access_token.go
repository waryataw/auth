package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (r repo) GetAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(r.authConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", fmt.Errorf("failed verifying refresh token: %w", err)
	}

	accessToken, err := utils.GenerateToken(models.User{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(r.authConfig.AccessTokenSecretKey()),
		time.Duration(r.authConfig.AccessTokenExpirationMinutes())*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("failed generating access token: %w", err)
	}

	return accessToken, nil
}
