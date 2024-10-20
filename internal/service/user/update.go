package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.User) error {
	if err := s.userRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
