package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/models"
)

// Create Метод создания пользователя.
func (s service) Create(ctx context.Context, user *models.User) (int64, error) {
	if !user.Role.IsValid() {
		return 0, fmt.Errorf("invalid user role")
	}

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}
