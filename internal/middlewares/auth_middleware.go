package middlewares

import (
	"net/http"
	"strings"

	"trithong.com/task-golang/internal/responses"
	"trithong.com/task-golang/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			responses.Error(c, http.StatusUnauthorized, "Missing authorization header")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			responses.Error(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

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
        c.Set("role", claims["role"].(string))

		c.Next()
	}
}
