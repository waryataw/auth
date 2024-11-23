package auth

import (
	"context"
	"fmt"
)

func (s service) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	newRefreshToken, err := s.repository.GetNewRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get new refresh token: %w", err)
	}

	return newRefreshToken, nil
}
