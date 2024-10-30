package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"github.com/waryataw/auth/internal/service/user"
	"github.com/waryataw/auth/internal/service/user/mocks"
)

func TestDelete(t *testing.T) {
	type mockBehavior func(mc *minimock.Controller) user.Repository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("failed to create user")

		id = gofakeit.Int64()
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
				req: id,
			},
			want: id,
			err:  nil,
			mockBehavior: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)

				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: 0,
			err:  fmt.Errorf("failed to delete user: %w", serviceErr),
			mockBehavior: func(mc *minimock.Controller) user.Repository {
				mock := mocks.NewRepositoryMock(mc)
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
			api := user.NewService(mock)

			err := api.Delete(tt.args.ctx, tt.args.req)

			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
