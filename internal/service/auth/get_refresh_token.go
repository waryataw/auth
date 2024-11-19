package auth

import "context"

func (s service) GetRefreshToken(_ context.Context, _ string) (string, error) {
	// TODO: Implement
	return "", nil
}
