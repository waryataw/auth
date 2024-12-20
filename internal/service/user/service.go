package user

import (
	"github.com/waryataw/auth/internal/api/user"
)

type service struct {
	repository Repository
}

// NewService Конструктор сервиса для операций с пользователем.
func NewService(
	repository Repository,
) user.Service {
	return &service{
		repository: repository,
	}
}
