package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (s service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepository.Get(ctx, 0, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrorInvalidCredentials
		}

		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if utils.VerifyPassword(user.Password, password) {
		token, err := s.repository.NewRefreshToken(user)
		if err != nil {
			return "", fmt.Errorf("failed to generate token: %w", err)
		}

		return token, nil
	}

	return "", models.ErrorInvalidCredentials
}
