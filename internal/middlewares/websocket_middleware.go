package middlewares

import (
	"net/http"

	"trithong.com/task-golang/internal/responses"
	"trithong.com/task-golang/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func WebSocketAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// WebSocket dùng query param thay vì header
		// Client kết nối: ws://localhost:8080/ws?token=xxx
		tokenString := c.Query("token")

		if tokenString == "" {
			responses.Error(c, http.StatusUnauthorized, "Missing token")
			c.Abort()
			return
		}

		token, err := utils.VerifyJWT(tokenString)

		if err != nil || !token.Valid {
			responses.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			responses.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)

		if !ok {
			responses.Error(c, http.StatusUnauthorized, "Invalid user id in token")
			c.Abort()
			return
		}

		userID := int(userIDFloat)

		c.Set("user_id", userID)

		c.Next()
	}
}