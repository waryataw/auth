package user

import (
	"github.com/waryataw/auth/pkg/userv1"
)

// Controller Имплементация Auth сервиса.
type Controller struct {
	userv1.UnimplementedUserServiceServer
	service MainService
}

// NewController Конструктор Имплементации Auth сервиса.
func NewController(service MainService) *Controller {
	return &Controller{
		service: service,
	}
}
