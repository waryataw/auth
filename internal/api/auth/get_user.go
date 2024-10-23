package auth

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/converter/auth"
	"github.com/waryataw/auth/pkg/authv1"
)

// GetUser Получение существующего пользователя
func (i *Implementation) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId(), req.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return auth.ToGetUserResponse(user), nil
}
