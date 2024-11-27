package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (r repo) GetRefreshToken(_ context.Context, user *models.User) (string, error) {
	token, err := utils.GenerateToken(
		*user,
		[]byte(r.authConfig.RefreshTokenSecretKey()),
		time.Duration(r.authConfig.RefreshTokenExpirationMinutes())*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
