package user

import (
	"github.com/waryataw/auth/internal/service/user"
	"github.com/waryataw/platform_common/pkg/db"
)

type repo struct {
	db db.Client
}

// NewRepository Конструктор репозитория пользователя.
func NewRepository(db db.Client) user.Repository {
	return &repo{db: db}
}
