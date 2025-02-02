package storage

import "go-auth/models"

type Cache interface {
	Get(key string) (*models.User, bool)
	Set(key string, value *models.User)
}

type RedisCache struct{}

func NewCache() (*RedisCache, error) {
	return &RedisCache{}, nil
}

func (r *RedisCache) Get(key string) (*models.User, bool) {
	return nil, false
}

func (r *RedisCache) Set(key string, value *models.User) {}
