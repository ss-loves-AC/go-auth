package handler

import (
	"go-auth/models"
	"go-auth/storage"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenHandler struct {
	DB    storage.DB
	Cache storage.Cache
}

func NewTokenHandler(db storage.DB, cache storage.Cache) *TokenHandler {
	return &TokenHandler{
		DB:    db,
		Cache: cache,
	}
}

func (h *TokenHandler) Authorize(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecretKey()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}

func (h *TokenHandler) RevokeToken(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecretKey()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token structure"})
		return
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token structure: missing JTI"})
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token structure: missing expiration"})
		return
	}

	expirationTime := time.Unix(int64(exp), 0)

	err = h.Cache.DeleteRefreshToken(c, jti)
	if err != nil {
		log.Println("Error deleting refresh token in Redis:", err)
	}

	h.Cache.SetRevokedToken(c, jti, &models.RevokedToken{
		UserID:    uint(userID),
		JTI:       jti,
		ExpiresAt: expirationTime,
	})

	err = h.DB.RevokeToken(c, jti)
	if err != nil {
		log.Println("Error storing revoked token in database:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}

func (h *TokenHandler) RefreshToken(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetSecretKey()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token structure"})
		return
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token structure"})
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token structure"})
		return
	}
	expirationTime := time.Unix(int64(exp), 0)

	if _, revoked := h.Cache.GetRevokedToken(c, jti); revoked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
		return
	}

	existingToken, err := h.DB.GetToken(c, jti)
	if err != nil || existingToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked or does not exist"})
		return
	}

	newToken, newJTI, err := GenerateNewJWTToken(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	h.Cache.SetRevokedToken(c, jti, &models.RevokedToken{
		UserID:    uint(userID),
		JTI:       jti,
		ExpiresAt: expirationTime,
	})

	if err := h.Cache.DeleteRefreshToken(c, jti); err != nil {
		log.Println("Error deleting old refresh token from cache:", err)
	}

	if err := h.DB.RevokeToken(c, jti); err != nil {
		log.Println("Error revoking token in DB:", err)
	}

	newRefreshToken := &models.RefreshToken{
		UserID:    uint(userID),
		JTI:       newJTI,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	h.Cache.SetRefreshToken(c, newJTI, newRefreshToken)

	if _, err := h.DB.CreateToken(c, newRefreshToken); err != nil {
		log.Println("Error storing new token in DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"new_token": newToken})
}

func GenerateNewJWTToken(userID uint) (string, string, error) {
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"user_id": userID,
		"jti":     jti,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(GetSecretKey()))
	if err != nil {
		return "", "", err
	}

	return tokenString, jti, nil
}

func GetSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set in environment variables")
	}
	return secretKey
}
