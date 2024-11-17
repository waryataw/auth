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
	"github.com/waryataw/auth/pkg/userv1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDeleteUser(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) user.MainService

	type args struct {
		ctx context.Context
		req *userv1.DeleteUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to delete user")

		id = gofakeit.Int64()

		req = &userv1.DeleteUserRequest{
			Id: id,
		}
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
			want: &emptypb.Empty{},
			err:  nil,
			mockBehavior: func(mc *minimock.Controller) user.MainService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)

				return mock
			},
		},
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  fmt.Errorf("failed to delete user: %w", serviceErr),
			mockBehavior: func(mc *minimock.Controller) user.MainService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)

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

			response, err := api.DeleteUser(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.want, response)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}

}
