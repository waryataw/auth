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

func TestGet(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) user.Repository

	type args struct {
		ctx  context.Context
		id   int64
		name string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to update user")

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
			ID:              id,
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            role,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               *models.User
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case by id",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: usr,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id, "").Return(usr, nil)

				return mock
			},
		},
		{
			name: "success case by name",
			args: args{
				ctx:  ctx,
				name: name,
			},
			want: usr,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.GetMock.Expect(ctx, 0, name).Return(usr, nil)

				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  fmt.Errorf("failed to get user: %w", serviceErr),
			userRepositoryMock: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id, "").Return(nil, serviceErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.userRepositoryMock(mc)
			api := user.NewService(mock)

			response, err := api.Get(tt.args.ctx, tt.args.id, tt.args.name)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
