package auth

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/converter/auth"
	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser Обновление существующего пользователя
func (i *Implementation) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, auth.ToUserForUpdate(req))
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
