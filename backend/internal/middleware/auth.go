package middleware

import (
	"strings"

	"github.com/codetesla51/todoapi/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization Header Required"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid Header Format"})
			c.Abort()
			return
		}
		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or Expired token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}

}
