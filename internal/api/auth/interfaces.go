package auth

import "context"

// Service Auth Сервис.
type Service interface {
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}
