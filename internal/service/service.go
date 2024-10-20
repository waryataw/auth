package service

import (
	"context"

	"github.com/waryataw/auth/internal/model"
)

// UserService Сервис для работы с Пользователем
type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}
