package access

import (
	"github.com/waryataw/auth/pkg/accessv1"
)

// Controller Имплементация Access сервиса.
type Controller struct {
	accessv1.UnimplementedAccessServer
	service Service
}

// NewController Конструктор Имплементации Auth сервиса.
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}
