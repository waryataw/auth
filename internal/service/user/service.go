package user

import (
	"github.com/waryataw/auth/internal/api/auth"
)

type service struct {
	repository Repository
}

// NewService Конструктор сервиса для операций с пользователем.
func NewService(
	repository Repository,
) auth.UserService {
	return &service{
		repository: repository,
	}
}
