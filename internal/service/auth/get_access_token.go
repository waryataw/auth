package auth

import (
	"context"
	"fmt"
)

func (s service) NewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	accessToken, err := s.repository.NewAccessToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	return accessToken, nil
}
