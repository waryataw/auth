package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/auth/internal/api/user"
	"github.com/waryataw/auth/internal/api/user/mocks"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/pkg/userv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpdateUser(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) user.MainService

	type args struct {
		ctx context.Context
		req *userv1.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to update user")

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		roles = []userv1.Role{
			userv1.Role_UNKNOWN,
			userv1.Role_USER,
			userv1.Role_ADMIN,
		}

		role = roles[gofakeit.Number(0, len(roles)-1)]

		userModel = &models.User{
			ID:    id,
			Name:  name,
			Email: email,
			Role:  models.Role(role),
		}

		req = &userv1.UpdateUserRequest{
			Id:    id,
			Name:  name,
			Email: email,
			Role:  role,
		}

		result = &emptypb.Empty{}
	)

	tests := []struct {
		name         string
		args         args
		want         *emptypb.Empty
		err          error
		mockBehavior mockBehavior
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: result,
			err:  nil,
			mockBehavior: func(mc *minimock.Controller) user.MainService {
				mock := mocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, userModel).Return(nil)

				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  fmt.Errorf("failed to update user: %w", serviceErr),
			mockBehavior: func(mc *minimock.Controller) user.MainService {
				mock := mocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, userModel).Return(serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.mockBehavior(mc)
			api := user.NewController(mock)

			response, err := api.UpdateUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
