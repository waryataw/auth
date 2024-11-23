package auth

import (
	"github.com/waryataw/auth/internal/api/auth"
)

type service struct {
	repository     Repository
	userRepository UserRepository
}

// NewService Конструктор сервиса Аутентификации.
func NewService(repository Repository, userRepository UserRepository) auth.Service {
	return &service{repository: repository, userRepository: userRepository}
}
