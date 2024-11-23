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
	GetRefreshToken(ctx context.Context, user *models.User) (string, error)
	GetNewRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}
