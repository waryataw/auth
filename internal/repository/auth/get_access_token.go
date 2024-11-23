package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (repository repo) GetAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(repository.authConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", fmt.Errorf("failed verifying refresh token: %w", err)
	}

	accessToken, err := utils.GenerateToken(models.User{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(repository.authConfig.AccessTokenSecretKey()),
		time.Duration(repository.authConfig.AccessTokenExpirationMinutes())*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("failed generating access token: %w", err)
	}

	return accessToken, nil
}
