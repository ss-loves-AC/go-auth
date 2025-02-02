package router

import (
	"go-auth/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *handler.AuthHandler, tokenHandler *handler.TokenHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/signup", authHandler.SignUp)

	r.POST("/signin", authHandler.Signin)

	r.POST("/authorize", tokenHandler.Authorize)

	r.POST("/revoke", tokenHandler.RevokeToken)

	r.POST("/refresh-token", tokenHandler.RefreshToken)

	return r
}
