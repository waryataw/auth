package auth

import (
	"github.com/waryataw/auth/pkg/authv1"
)

// Controller Имплементация Auth сервиса
type Controller struct {
	authv1.UnimplementedAuthServiceServer
	userService UserService
}

// NewController Конструктор Имплементации Auth сервиса
func NewController(userService UserService) *Controller {
	return &Controller{
		userService: userService,
	}
}
