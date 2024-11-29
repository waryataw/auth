package auth

import (
	"context"

	"github.com/waryataw/auth/internal/models"
)

// UserRepository Интерфейс репозитория для операций с пользователем.
type UserRepository interface {
	Get(ctx context.Context, id int64, name string) (*models.User, error)
}

// Repository Интерфейс Auth репозитория.
type Repository interface {
	NewRefreshToken(user *models.User) (string, error)
	UpdateRefreshToken(oldRefreshToken string) (string, error)
	NewAccessToken(refreshToken string) (string, error)
}
