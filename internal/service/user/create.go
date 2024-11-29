package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

// Create Метод создания пользователя.
func (s service) Create(ctx context.Context, user *models.User) (int64, error) {
	if !user.Role.IsValid() {
		return 0, fmt.Errorf("invalid user role")
	}

	passwordHashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = passwordHashed

	confirmationHashed, err := utils.HashPassword(user.PasswordConfirm)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password confirm: %w", err)
	}

	user.PasswordConfirm = confirmationHashed

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}
