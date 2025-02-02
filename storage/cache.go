package storage

import (
	"context"
	"go-auth/models"
	"log"

	"github.com/go-redis/redis/v8"
)

type Cache interface {

	GetUser(ctx context.Context, key string) (*models.User, bool)    // Retrieve user data by key
	SetUser(ctx context.Context, key string, value *models.User)    // Store user data

	GetRefreshToken(ctx context.Context, jti string) (*models.RefreshToken, bool) // Retrieve refresh token by JTI
	SetRefreshToken(ctx context.Context, jti string, value *models.RefreshToken)   // Store refresh token

	GetRevokedToken(ctx context.Context, jti string) (*models.RevokedToken, bool) // Retrieve revoked token by JTI
	SetRevokedToken(ctx context.Context, jti string, value *models.RevokedToken)   // Store revoked token
}

type RedisCache struct {
	Client *redis.Client
}

func NewCache() (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	return &RedisCache{
		Client: client,
	}, nil
}

func (r *RedisCache) GetUser(ctx context.Context, key string) (*models.User, bool)

func (r *RedisCache) SetUser(ctx context.Context, key string, value *models.User)

func (r *RedisCache) GetRefreshToken(ctx context.Context, jti string) (*models.RefreshToken, bool)

func (r *RedisCache) SetRefreshToken(ctx context.Context, jti string, value *models.RefreshToken)

func (r *RedisCache) GetRevokedToken(ctx context.Context, jti string) (*models.RevokedToken, bool)

func (r *RedisCache) SetRevokedToken(ctx context.Context, jti string, value *models.RevokedToken)
