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
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) user.Service

	type args struct {
		ctx context.Context
		req *userv1.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to create user")

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(
			true,
			true,
			true,
			false,
			false,
			0,
		)

		roles = []userv1.Role{
			userv1.Role_UNKNOWN,
			userv1.Role_USER,
			userv1.Role_ADMIN,
		}

		role = roles[gofakeit.Number(0, len(roles)-1)]

		userModel = &models.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            models.Role(role),
		}

		req = &userv1.CreateUserRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
		}

		res = &userv1.CreateUserResponse{
			Id: id,
		}
	)

	tests := []struct {
		name         string
		args         args
		want         *userv1.CreateUserResponse
		err          error
		mockBehavior mockBehavior
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			mockBehavior: func(mc *minimock.Controller) user.Service {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userModel).Return(id, nil)

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
			err:  fmt.Errorf("failed to create user: %w", serviceErr),
			mockBehavior: func(mc *minimock.Controller) user.Service {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userModel).Return(0, serviceErr)

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

			response, err := api.CreateUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
