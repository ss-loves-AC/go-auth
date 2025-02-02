package storage

import (
	"context"
	"fmt"
	"go-auth/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbHost     = "mysql"
	dbName     = "go_auth"
	dbUser     = "go_user"
	dbPassword = "go_password"
	dbPort     = "3306"
)

type MysqlDB struct {
	DB *gorm.DB
}

func NewDB() (*MysqlDB, error) {

	source := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
	if err != nil {
		log.Println("gorm Db connection ", err)
		return nil, fmt.Errorf("gorm DB connection error: %w", err)
	}

	tables := []interface{}{
		&models.User{},
		&models.RefreshToken{},
		&models.RevokedToken{},
	}

	if err := db.AutoMigrate(tables...); err != nil {
		fmt.Println("Failed to auto migrate:", err)
		return nil, fmt.Errorf("failed to auto migrate: %w", err)

	}

	return &MysqlDB{DB: db}, nil
}

func (db *MysqlDB) GetUser(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := db.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (db *MysqlDB) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := db.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (db *MysqlDB) GetToken(ctx context.Context, jti string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	if err := db.DB.WithContext(ctx).Where("jti = ?", jti).First(&token).Error; err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return &token, nil
}

func (db *MysqlDB) CreateToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error) {
	if err := db.DB.WithContext(ctx).Create(token).Error; err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	return token, nil
}

func (db *MysqlDB) RevokeToken(ctx context.Context, jti string) error {

	if err := db.DB.WithContext(ctx).Where("jti = ?", jti).Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	var token models.RefreshToken
	err := db.DB.WithContext(ctx).Where("jti = ?", jti).First(&token).Error
	if err == nil {
		revokedToken := &models.RevokedToken{
			UserID:    token.UserID,
			JTI:       token.JTI,
			ExpiresAt: token.ExpiresAt,
		}
		err = db.DB.WithContext(ctx).Create(&revokedToken).Error
		if err != nil {
			return fmt.Errorf("failed to create revoked token: %w", err)
		}
	}

	return nil
}

func (db *MysqlDB) RefreshToken(ctx context.Context, jti string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	if err := db.DB.WithContext(ctx).Where("jti = ?", jti).First(&token).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return &token, nil
}
