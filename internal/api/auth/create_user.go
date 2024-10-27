package auth

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/converter/auth"
	"github.com/waryataw/auth/pkg/authv1"
)

// CreateUser Добавление нового пользователя
func (c Controller) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	id, err := c.userService.Create(ctx, auth.ToUser(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &authv1.CreateUserResponse{Id: id}, nil
}
