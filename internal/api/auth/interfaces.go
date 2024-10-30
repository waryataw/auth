package auth

import (
	"context"

	"github.com/waryataw/auth/internal/models"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i github.com/waryataw/auth/internal/api/auth.* -o "./mocks/mocks.go"

// UserService Сервис для работы с Пользователем.
type UserService interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	Get(ctx context.Context, id int64, name string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}
