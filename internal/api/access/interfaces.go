package access

import (
	"context"
)

// Service Access Сервис.
type Service interface {
	Check(ctx context.Context, endpointAddress string) error
}
