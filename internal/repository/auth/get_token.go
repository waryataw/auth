package auth

import (
	"fmt"
	"time"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (repository repo) GetToken(user *models.User) (string, error) {
	token, err := utils.GenerateToken(
		*user,
		[]byte(repository.refreshTokenConfig.RefreshTokenSecretKey()),
		time.Duration(repository.refreshTokenConfig.RefreshTokenExpirationMinutes())*time.Minute,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
