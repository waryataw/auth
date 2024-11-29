package access

import (
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/service/access"
	"github.com/waryataw/platform_common/pkg/db"
)

type repo struct {
	authConfig config.AuthConfig
	db         db.Client
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(authConfig config.AuthConfig, db db.Client) access.Repository {
	return &repo{authConfig: authConfig, db: db}
}
