package auth

import (
	"github.com/waryataw/auth/pkg/authv1"
)

// Controller Имплементация Auth сервиса.
type Controller struct {
	authv1.UnimplementedAuthServer
	service Service
}

// NewController Конструктор Имплементации Auth сервиса.
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}
