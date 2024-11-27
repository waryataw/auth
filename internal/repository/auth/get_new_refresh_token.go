package auth

import (
	"fmt"
	"time"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (r repo) UpdateRefreshToken(oldRefreshToken string) (string, error) {
	claims, err := utils.VerifyToken(oldRefreshToken, []byte(r.authConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", fmt.Errorf("verifying old refresh token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(models.User{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(r.authConfig.RefreshTokenSecretKey()),
		time.Duration(r.authConfig.RefreshTokenExpirationMinutes())*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("generating refresh token: %w", err)
	}

	return refreshToken, nil
}
