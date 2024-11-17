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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUser(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) user.Service

	type args struct {
		ctx context.Context
		req *userv1.GetUserRequest
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

		roles = []userv1.Role{
			userv1.Role_UNKNOWN,
			userv1.Role_USER,
			userv1.Role_ADMIN,
		}

		role = roles[gofakeit.Number(0, len(roles)-1)]

		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		userModel = &models.User{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            models.Role(role),
			CreatedAt:       &createdAt,
			UpdatedAt:       &updatedAt,
		}

		reqByID = &userv1.GetUserRequest{
			Query: &userv1.GetUserRequest_Id{
				Id: id,
			},
		}

		reqByName = &userv1.GetUserRequest{
			Query: &userv1.GetUserRequest_Name{
				Name: name,
			},
		}

		res = &userv1.GetUserResponse{
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: timestamppb.New(*userModel.CreatedAt),
			UpdatedAt: timestamppb.New(*userModel.UpdatedAt),
		}
	)

	tests := []struct {
		name         string
		args         args
		want         *userv1.GetUserResponse
		err          error
		mockBehavior mockBehavior
	}{
		{
			name: "success case by id",
			args: args{
				ctx: ctx,
				req: reqByID,
			},
			want: res,
			err:  nil,
			mockBehavior: func(mc *minimock.Controller) user.Service {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id, "").Return(userModel, nil)

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
			mockBehavior: func(mc *minimock.Controller) user.Service {
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
			mockBehavior: func(mc *minimock.Controller) user.Service {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, 0, name).Return(userModel, nil)

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
			mockBehavior: func(mc *minimock.Controller) user.Service {
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

			mock := tt.mockBehavior(mc)
			api := user.NewController(mock)

			response, err := api.GetUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}

}
