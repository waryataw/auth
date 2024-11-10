package consumer

import (
	"context"
)

// Service Сервис слушателя.
type Service interface {
	RunConsumer(ctx context.Context) error
}
