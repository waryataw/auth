package user

import (
	"github.com/waryataw/auth/internal/repository"
	"github.com/waryataw/auth/internal/service"
)

type userService struct {
	userRepository repository.UserRepository
}

// NewService Конструктор сервиса для операций с пользователем
func NewService(
	userRepository repository.UserRepository,
) service.UserService {
	return &userService{
		userRepository: userRepository,
	}
}
