package auth

import (
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/service/auth"
)

type repo struct {
	refreshTokenConfig config.RefreshTokenConfig
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(refreshTokenConfig config.RefreshTokenConfig) auth.Repository {
	return &repo{refreshTokenConfig: refreshTokenConfig}
}
