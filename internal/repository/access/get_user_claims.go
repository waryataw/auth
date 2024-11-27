package access

import (
	"context"
	"fmt"

	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/internal/utils"
)

func (r repo) GetUserClaims(_ context.Context, accessToken string) (*models.UserClaims, error) {
	claims, err := utils.VerifyToken(accessToken, []byte(r.authConfig.AccessTokenSecretKey()))
	if err != nil {
		return nil, fmt.Errorf("access token verification failed: %w", err)
	}

	return claims, nil
}
