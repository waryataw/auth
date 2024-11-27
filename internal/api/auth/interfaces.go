package auth

import "context"

// Service Auth Сервис.
type Service interface {
	NewAccessToken(ctx context.Context, refreshToken string) (string, error)
	UpdateRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}
