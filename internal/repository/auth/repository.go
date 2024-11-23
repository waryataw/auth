package auth

import (
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/service/auth"
)

type repo struct {
	authConfig config.AuthConfig
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(authConfig config.AuthConfig) auth.Repository {
	return &repo{authConfig: authConfig}
}
