package auth

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser Удаление существующего пользователя
func (i *Implementation) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
