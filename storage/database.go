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


func (db *MysqlDB) GetUser(ctx context.Context, email string) (*models.User, error)

func (db *MysqlDB) CreateUser(ctx context.Context, user *models.User) (*models.User, error)

func (db *MysqlDB) GetToken(ctx context.Context, jti string) (*models.RefreshToken, error)

func (db *MysqlDB) CreateToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error)

func (db *MysqlDB) RevokeToken(ctx context.Context, jti string) error

func (db *MysqlDB) RefreshToken(ctx context.Context, jti string) (*models.RefreshToken, error)