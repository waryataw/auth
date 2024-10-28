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
)

func TestCreateUser(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) auth.UserService

	type args struct {
		ctx context.Context
		req *authv1.CreateUserRequest
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

		roles = []authv1.Role{
			authv1.Role_UNKNOWN,
			authv1.Role_USER,
			authv1.Role_ADMIN,
		}

		role = roles[gofakeit.Number(0, len(roles)-1)]

		user = &models.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            models.Role(role),
		}

		req = &authv1.CreateUserRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
		}

		res = &authv1.CreateUserResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *authv1.CreateUserResponse
		err             error
		noteServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
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
			noteServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.noteServiceMock(mc)
			api := auth.NewController(mock)

			response, err := api.CreateUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
