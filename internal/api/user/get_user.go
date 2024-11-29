package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/converter/auth"
	"github.com/waryataw/auth/pkg/userv1"
)

// GetUser Получение существующего пользователя.
func (c Controller) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := c.service.Get(ctx, req.GetId(), req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return auth.ToGetUserResponse(user), nil
}
