package kafka

import (
	"context"

	"github.com/waryataw/auth/internal/client/kafka/consumer"
)

// Consumer Слушатель.
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
