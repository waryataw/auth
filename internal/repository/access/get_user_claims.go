package access

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (repository repo) GetUserClaims(_ context.Context, accessToken string) (*models.UserClaims, error) {
	claims, err := utils.VerifyToken(accessToken, []byte(repository.authConfig.AccessTokenSecretKey()))
	if err != nil {
		return nil, fmt.Errorf("access token verification failed: %w", err)
	}

	return claims, nil
}
