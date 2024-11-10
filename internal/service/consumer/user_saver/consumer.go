package user_saver

import (
	"context"

	"github.com/waryataw/auth/internal/client/kafka"
	def "github.com/waryataw/auth/internal/service/consumer"
	"github.com/waryataw/auth/internal/service/user"
)

var _ def.Service = (*service)(nil)

type service struct {
	userRepository user.Repository
	consumer       kafka.Consumer
}

// NewService Конструктор сервиса слушателя.
func NewService(
	userRepository user.Repository,
	consumer kafka.Consumer,
) *service {
	return &service{
		userRepository: userRepository,
		consumer:       consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, "user-topic", s.UserSaveHandler)
	}()

	return errChan
}
