package auth

import (
	"context"
	"fmt"
)

func (s service) UpdateRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	newRefreshToken, err := s.repository.UpdateRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get new refresh token: %w", err)
	}

	return newRefreshToken, nil
}
