package storage

import (
	"context"
	"go-auth/models"
)

type Cache interface {
	GetUser(ctx context.Context, key string) (*models.User, bool) // Retrieve user data by key
	SetUser(ctx context.Context, key string, value *models.User)  // Store user data

	GetRefreshToken(ctx context.Context, jti string) (*models.RefreshToken, bool) // Retrieve refresh token by JTI
	SetRefreshToken(ctx context.Context, jti string, value *models.RefreshToken) (error) // Store refresh token
	DeleteRefreshToken(ctx context.Context, jti string) error // delete refresh token

	GetRevokedToken(ctx context.Context, jti string) (*models.RevokedToken, bool) // Retrieve revoked token by JTI
	SetRevokedToken(ctx context.Context, jti string, value *models.RevokedToken)  // Store revoked token
}
