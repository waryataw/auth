package access

import (
	"context"
	"fmt"
)

func (s service) Check(ctx context.Context, accessToken string, endpointAddress string) error {
	claims, err := s.repository.GetUserClaims(accessToken)
	if err != nil {
		return fmt.Errorf("failed to get user claims: %w", err)
	}

	accessibleRolesMap, err := s.repository.GetAccessibleRoles(ctx, endpointAddress)
	if err != nil {
		return fmt.Errorf("failed to get accessible roles: %w", err)
	}

	role, ok := accessibleRolesMap[endpointAddress]
	if !ok {
		return nil
	}

	if role == claims.Role {
		return nil
	}

	return fmt.Errorf("access to endpoint %s is not authorized to access this resource", endpointAddress)
}
