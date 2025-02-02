package handler

import (
	"go-auth/models"
	"go-auth/storage"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB    storage.DB
	Cache storage.Cache
}

func NewAuthHandler(db storage.DB, cache storage.Cache) *AuthHandler {
	return &AuthHandler{
		DB:    db,
		Cache: cache,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !isValidEmail(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	_, err := h.DB.GetUser(c, input.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	user := &models.User{
		Email:        input.Email,
		HashPassword: string(hashedPassword),
	}
	_, err = h.DB.CreateUser(c, user)
	if err != nil {
		log.Println("Error creating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.DB.GetUser(c, input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, jti, err := GenerateNewJWTToken(uint(user.ID)) // Convert to uint here
	if err != nil {
		log.Println("Error generating JWT token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken := &models.RefreshToken{
		UserID:    uint(user.ID), // Convert to uint here
		JTI:       jti,
		ExpiresAt: time.Now().Add(time.Hour * 24), // Set expiration to 24 hours
	}

	if err := h.Cache.SetRefreshToken(c, jti, refreshToken); err != nil {
		log.Println("Error storing refresh token in cache:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token in cache"})
		return
	}

	if _, err := h.DB.CreateToken(c, refreshToken); err != nil {
		log.Println("Error storing refresh token in database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken.JTI,
	})
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}
