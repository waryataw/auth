package access

import (
	"github.com/waryataw/auth/internal/api/access"
)

type service struct {
}

// NewService Конструктор Сервиса контроля доступа.
func NewService() access.Service {
	return &service{}
}
