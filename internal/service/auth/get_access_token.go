package auth

import (
	"fmt"
)

func (s service) NewAccessToken(refreshToken string) (string, error) {
	accessToken, err := s.repository.NewAccessToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	return accessToken, nil
}
