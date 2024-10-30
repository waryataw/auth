package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/service/user"
	"github.com/waryataw/auth/internal/service/user/mocks"
)

func TestCreate(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) user.Repository

	type args struct {
		ctx context.Context
		req *models.User
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

		roles = []models.Role{
			models.RoleUnknown,
			models.RoleUser,
			models.RoleAdmin,
		}

		role = roles[gofakeit.Number(0, len(roles)-1)]

		usr = &models.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
		}

		usrInvalidRole = &models.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            models.Role(666),
		}
	)

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		mockBehavior mockBehavior
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: usr,
			},
			want: id,
			err:  nil,
			mockBehavior: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, usr).Return(id, nil)

				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: usr,
			},
			want: 0,
			err:  fmt.Errorf("failed to create user: %w", serviceErr),
			mockBehavior: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, usr).Return(0, serviceErr)

				return mock
			},
		},
		{
			name: "service error invalid role case",
			args: args{
				ctx: ctx,
				req: usrInvalidRole,
			},
			want: 0,
			err:  fmt.Errorf("invalid user role"),
			mockBehavior: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.mockBehavior(mc)
			api := user.NewService(mock)

			response, err := api.Create(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
