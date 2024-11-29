package user

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/pkg/userv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser Удаление существующего пользователя.
func (c Controller) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*emptypb.Empty, error) {
	err := c.service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
