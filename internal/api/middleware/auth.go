package middleware

import (
	"net/http"
	"strings"

	"todo-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果是 WebSocket 请求，跳过 HTTP 认证
		if strings.Contains(c.Request.Header.Get("Upgrade"), "websocket") {
			c.Next()
			return
		}

		// 其他请求（HTTP）需要认证
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 从 Bearer token 中提取 Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 验证 Token
		userID, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户 ID 存储到上下文中
		c.Set("userID", userID)
		c.Next()
	}
}
