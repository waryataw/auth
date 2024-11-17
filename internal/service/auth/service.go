package auth

import (
	"github.com/waryataw/auth/internal/api/auth"
)

type service struct {
}

// NewService Конструктор сервиса Аутентификации.
func NewService() auth.Service {
	return &service{}
}
