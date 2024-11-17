package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/converter/auth"
	"github.com/waryataw/auth/pkg/userv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser Обновление существующего пользователя.
func (c Controller) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*emptypb.Empty, error) {
	err := c.service.Update(ctx, auth.ToUserForUpdate(req))
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
