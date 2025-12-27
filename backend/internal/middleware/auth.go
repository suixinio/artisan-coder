package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"artisan-coder/pkg/jwt"
	"artisan-coder/pkg/response"
)

const (
	userIDKey = "user_id"
)

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization token")
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// 将用户 ID 存储到上下文（存储为字符串）
		c.Set(userIDKey, claims.UserID.String())
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}
