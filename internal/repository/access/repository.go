package access

import (
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/service/access"
)

type repo struct {
	authConfig config.AuthConfig
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(authConfig config.AuthConfig) access.Repository {
	return &repo{authConfig: authConfig}
}
