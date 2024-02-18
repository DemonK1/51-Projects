package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_gin_jwt/helpers"
	"net/http"
)

func Authenticate(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("没有提供token")})
		c.Abort()
		return
	}

	claims, err := helpers.ValidateToken(clientToken)
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Set("first_name", claims.FirstName)
	c.Set("last_name", claims.LastName)
	c.Set("uid", claims.Uid)
	c.Set("user_type", claims.UserType)
	c.Next()
}

// func Authenticate(context *gin.Context) {
//
// }
