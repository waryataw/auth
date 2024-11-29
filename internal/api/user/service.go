package user

import (
	"github.com/waryataw/auth/pkg/userv1"
)

// Controller Имплементация User сервиса.
type Controller struct {
	userv1.UnimplementedUserServiceServer
	service Service
}

// NewController Конструктор Имплементации User сервиса.
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}
