package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/converter/auth"
	"github.com/waryataw/auth/pkg/userv1"
)

// CreateUser Добавление нового пользователя.
func (c Controller) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	id, err := c.service.Create(ctx, auth.ToUser(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &userv1.CreateUserResponse{Id: id}, nil
}
