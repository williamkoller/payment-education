package auth_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	port_auth_cryptography "github.com/williamkoller/system-education/internal/auth/port/cryptography"
)

func AuthMiddleware(jwt port_auth_cryptography.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		tokenStr := authHeader[7:]
		claims, err := jwt.Verify(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("userEmail", claims["email"])

		if userID, ok := claims["user_id"]; ok {
			c.Set("userID", userID)
		}

		if modules, ok := claims["modules"]; ok {
			c.Set("modules", modules)
		}

		if actions, ok := claims["actions"]; ok {
			c.Set("actions", actions)
		}

		c.Next()
	}
}
