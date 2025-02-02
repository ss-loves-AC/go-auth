package storage

import (
	"context"
	"go-auth/models"

	"gorm.io/gorm"
)

type DB interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)

	GetToken(ctx context.Context, jti string) (*models.RefreshToken, error)
	CreateToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, jti string) error
	RefreshToken(ctx context.Context, jti string) (*models.RefreshToken, error)
}

type MysqlDB struct {
	DB *gorm.DB
}
