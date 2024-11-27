package auth

import (
	"fmt"
)

func (s service) UpdateRefreshToken(oldRefreshToken string) (string, error) {
	newRefreshToken, err := s.repository.UpdateRefreshToken(oldRefreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get new refresh token: %w", err)
	}

	return newRefreshToken, nil
}
