package user

import (
	"context"

	"github.com/waryataw/auth/internal/models"
)

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../../bin/minimock -i github.com/waryataw/user/internal/api/user.* -o "./mocks/mocks.go"

// Service Сервис для работы с Пользователем.
type Service interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	Get(ctx context.Context, id int64, name string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}
