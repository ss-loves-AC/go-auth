package storage

import (
	"context"
	"encoding/json"
	"go-auth/models"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func (r *RedisCache) GetUser(ctx context.Context, key string) (*models.User, bool) {
	val, err := r.Client.Get(ctx, "user:"+key).Result()
	if err == redis.Nil {
		log.Println("User not found in Redis")
		return nil, false
	} else if err != nil {
		log.Println("Error retrieving user from Redis:", err)
		return nil, false
	}

	var user models.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		log.Println("Error unmarshalling user from Redis:", err)
		return nil, false
	}

	return &user, true
}

func (r *RedisCache) SetUser(ctx context.Context, key string, value *models.User) {
	userJSON, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshalling user for Redis:", err)
		return
	}

	err = r.Client.Set(ctx, "user:"+key, userJSON, time.Hour*24).Err()
	if err != nil {
		log.Println("Error setting user in Redis:", err)
	}
}

func (r *RedisCache) GetRefreshToken(ctx context.Context, jti string) (*models.RefreshToken, bool) {
	val, err := r.Client.Get(ctx, "refresh_token:"+jti).Result()
	if err == redis.Nil {
		log.Println("Refresh token not found in Redis")
		return nil, false
	} else if err != nil {
		log.Println("Error retrieving refresh token from Redis:", err)
		return nil, false
	}

	var token models.RefreshToken
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		log.Println("Error unmarshalling refresh token from Redis:", err)
		return nil, false
	}

	return &token, true
}

func (r *RedisCache) SetRefreshToken(ctx context.Context, jti string, value *models.RefreshToken) {
	tokenJSON, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshalling refresh token for Redis:", err)
		return
	}

	err = r.Client.Set(ctx, "refresh_token:"+jti, tokenJSON, time.Until(value.ExpiresAt)).Err()
	if err != nil {
		log.Println("Error setting refresh token in Redis:", err)
	}
}

func (r *RedisCache) GetRevokedToken(ctx context.Context, jti string) (*models.RevokedToken, bool) {
	val, err := r.Client.Get(ctx, "revoked_token:"+jti).Result()
	if err == redis.Nil {
		log.Println("Revoked token not found in Redis")
		return nil, false
	} else if err != nil {
		log.Println("Error retrieving revoked token from Redis:", err)
		return nil, false
	}

	var token models.RevokedToken
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		log.Println("Error unmarshalling revoked token from Redis:", err)
		return nil, false
	}

	return &token, true
}

func (r *RedisCache) SetRevokedToken(ctx context.Context, jti string, value *models.RevokedToken) {
	tokenJSON, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshalling revoked token for Redis:", err)
		return
	}

	err = r.Client.Set(ctx, "revoked_token:"+jti, tokenJSON, time.Until(value.ExpiresAt)).Err()
	if err != nil {
		log.Println("Error setting revoked token in Redis:", err)
	}
}
