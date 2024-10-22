package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/model"
)

func (s *userService) Get(ctx context.Context, id int64, name string) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
