package user

import (
	"github.com/waryataw/auth/internal/client/db"
	"github.com/waryataw/auth/internal/repository"
)

type repo struct {
	db db.Client
}

// NewRepository Конструктор репозитория пользователя
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
