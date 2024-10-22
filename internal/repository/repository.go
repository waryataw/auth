package repository

import (
	"context"

	"github.com/waryataw/auth/internal/model"
)

// UserRepository Интерфейс репозитория для операций с пользователем
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64, name string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}
