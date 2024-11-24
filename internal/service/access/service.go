package access

import (
	"github.com/waryataw/auth/internal/api/access"
)

type service struct {
	repository Repository
}

// NewService Конструктор Сервиса контроля доступа.
func NewService(repository Repository) access.Service {
	return &service{repository: repository}
}
