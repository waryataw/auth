package access

import (
	"context"
	"strings"

	"github.com/waryataw/auth/pkg/accessv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const authPrefix = "Bearer "

// Check Проверка доступа.
func (c Controller) Check(ctx context.Context, req *accessv1.CheckRequest) (*emptypb.Empty, error) {
	if req == nil || req.EndpointAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "endpoint address is required")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata is required")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.InvalidArgument, "authorization header is required")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, status.Error(codes.InvalidArgument, "invalid authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	if err := c.service.Check(ctx, accessToken, req.EndpointAddress); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &emptypb.Empty{}, nil
}
