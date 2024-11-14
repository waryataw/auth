package user_saver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/waryataw/auth/internal/models"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &models.User{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return fmt.Errorf("failed to parsing user: %w", err)
	}

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("User with id %d created\n", id)

	return nil
}
