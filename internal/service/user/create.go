package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/model"
)

func (s *userService) Create(ctx context.Context, user *model.User) (int64, error) {
	if !user.Role.IsValid() {
		return 0, fmt.Errorf("user role is not valid")
	}

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}
