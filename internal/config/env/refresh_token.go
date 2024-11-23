package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/waryataw/auth/internal/config"
)

type refreshTokenConfig struct {
	refreshTokenSecretKey         string
	refreshTokenExpirationMinutes int64
}

// NewRefreshTokenConfig Конструктор конфига Refresh токена.
func NewRefreshTokenConfig() (config.RefreshTokenConfig, error) {
	refreshTokenSecretKey := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if len(refreshTokenSecretKey) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", "REFRESH_TOKEN_SECRET_KEY")
	}

	refreshTokenExpirationMinutes, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRATION_MINUTES"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s is not set", "REFRESH_TOKEN_EXPIRATION_MINUTES")
	}

	return &refreshTokenConfig{
		refreshTokenSecretKey:         refreshTokenSecretKey,
		refreshTokenExpirationMinutes: refreshTokenExpirationMinutes,
	}, nil
}

func (r refreshTokenConfig) RefreshTokenSecretKey() string {
	return r.refreshTokenSecretKey
}

func (r refreshTokenConfig) RefreshTokenExpirationMinutes() int64 {
	return r.refreshTokenExpirationMinutes
}
