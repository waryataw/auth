package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/waryataw/auth/internal/config"
)

type authConfig struct {
	authPrefix                    string
	refreshTokenSecretKey         string
	refreshTokenExpirationMinutes int64
	accessTokenSecretKey          string
	accessTokenExpirationMinutes  int64
}

// NewAuthConfig Конструктор конфига Refresh токена.
func NewAuthConfig() (config.AuthConfig, error) {
	authPrefix := os.Getenv("AUTH_PREFIX")
	if len(authPrefix) == 0 {
		return nil, fmt.Errorf("environment variable AUTH_PREFIX is not set")
	}

	refreshTokenSecretKey := os.Getenv("REFRESH_TOKEN_SECRET_KEY")
	if len(refreshTokenSecretKey) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", "REFRESH_TOKEN_SECRET_KEY")
	}

	refreshTokenExpirationMinutes, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRATION_MINUTES"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s is not set", "REFRESH_TOKEN_EXPIRATION_MINUTES")
	}

	accessTokenSecretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if len(accessTokenSecretKey) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", "ACCESS_TOKEN_SECRET_KEY")
	}

	accessTokenExpirationMinutes, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRATION_MINUTES"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s is not set", "ACCESS_TOKEN_EXPIRATION_MINUTES")
	}

	return &authConfig{
		authPrefix:                    authPrefix,
		refreshTokenSecretKey:         refreshTokenSecretKey,
		refreshTokenExpirationMinutes: refreshTokenExpirationMinutes,
		accessTokenSecretKey:          accessTokenSecretKey,
		accessTokenExpirationMinutes:  accessTokenExpirationMinutes,
	}, nil
}

func (r authConfig) RefreshTokenSecretKey() string {
	return r.refreshTokenSecretKey
}

func (r authConfig) RefreshTokenExpirationMinutes() int64 {
	return r.refreshTokenExpirationMinutes
}

func (r authConfig) AuthPrefix() string {
	return r.authPrefix
}

func (r authConfig) AccessTokenSecretKey() string {
	return r.accessTokenSecretKey
}

func (r authConfig) AccessTokenExpirationMinutes() int64 {
	return r.accessTokenExpirationMinutes
}
