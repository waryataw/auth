package auth

import (
	"github.com/waryataw/auth/internal/service"
	"github.com/waryataw/auth/pkg/authv1"
)

// Implementation Имплементация Auth сервиса
type Implementation struct {
	authv1.UnimplementedAuthServiceServer
	userService service.UserService
}

// NewImplementation Конструктор Имплементации Auth сервиса
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
