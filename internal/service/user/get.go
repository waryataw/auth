package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/models"
)

// Get Метод получения пользователя.
func (s service) Get(ctx context.Context, id int64, name string) (*models.User, error) {
	user, err := s.repository.Get(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
