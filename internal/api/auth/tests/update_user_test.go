package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/auth/internal/api/auth"
	"github.com/waryataw/auth/internal/api/auth/mocks"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/pkg/authv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpdateUser(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) auth.UserService

	type args struct {
		ctx context.Context
		req *authv1.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to update user")

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		roles = []authv1.Role{
			authv1.Role_UNKNOWN,
			authv1.Role_USER,
			authv1.Role_ADMIN,
		}

		role = roles[gofakeit.Number(0, len(roles)-1)]

		user = &models.User{
			ID:    id,
			Name:  name,
			Email: email,
			Role:  models.Role(role),
		}

		req = &authv1.UpdateUserRequest{
			Id:    id,
			Name:  name,
			Email: email,
			Role:  role,
		}

		result = &emptypb.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: result,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, user).Return(nil)

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
			userServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, user).Return(serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.userServiceMock(mc)
			api := auth.NewController(mock)

			response, err := api.UpdateUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
