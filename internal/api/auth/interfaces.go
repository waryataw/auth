package auth

import "context"

// Service Auth Сервис.
type Service interface {
	NewAccessToken(refreshToken string) (string, error)
	UpdateRefreshToken(oldRefreshToken string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}
