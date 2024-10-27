package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/models"
)

// Update Метод изменения пользователя
func (s service) Update(ctx context.Context, user *models.User) error {
	if !user.Role.IsValid() {
		return fmt.Errorf("invalid user role")
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
