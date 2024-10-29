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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUser(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) auth.UserService

	type args struct {
		ctx context.Context
		req *authv1.GetUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to get user")

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

		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		user = &models.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            models.Role(role),
			CreatedAt:       &createdAt,
			UpdatedAt:       &updatedAt,
		}

		reqByID = &authv1.GetUserRequest{
			Query: &authv1.GetUserRequest_Id{
				Id: id,
			},
		}

		reqByName = &authv1.GetUserRequest{
			Query: &authv1.GetUserRequest_Name{
				Name: name,
			},
		}

		res = &authv1.GetUserResponse{
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: timestamppb.New(*user.CreatedAt),
			UpdatedAt: timestamppb.New(*user.UpdatedAt),
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *authv1.GetUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case by id",
			args: args{
				ctx: ctx,
				req: reqByID,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id, "").Return(user, nil)

				return mock
			},
		},
		{
			name: "service error case by id",
			args: args{
				ctx: ctx,
				req: reqByID,
			},
			want: nil,
			err:  fmt.Errorf("failed to get user: %w", serviceErr),
			userServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id, "").Return(nil, serviceErr)

				return mock
			},
		},
		{
			name: "success case by name",
			args: args{
				ctx: ctx,
				req: reqByName,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, 0, name).Return(user, nil)

				return mock
			},
		},
		{
			name: "service error case by name",
			args: args{
				ctx: ctx,
				req: reqByName,
			},
			want: nil,
			err:  fmt.Errorf("failed to get user: %w", serviceErr),
			userServiceMock: func(mc *minimock.Controller) auth.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, 0, name).Return(nil, serviceErr)

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

			response, err := api.GetUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}

}
