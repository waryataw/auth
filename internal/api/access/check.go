package access

import (
	"context"

	"github.com/waryataw/auth/pkg/accessv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Check Проверка доступа.
func (c Controller) Check(ctx context.Context, req *accessv1.CheckRequest) (*emptypb.Empty, error) {
	if req == nil || req.EndpointAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "endpoint address is required")
	}

	if err := c.service.Check(ctx, req.EndpointAddress); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &emptypb.Empty{}, nil
}
