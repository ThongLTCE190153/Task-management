package middlewares

import (
	"net/http"

	"trithong.com/task-golang/internal/responses"

	"github.com/gin-gonic/gin"
)

// Dùng như này: RoleMiddleware("admin") hoặc RoleMiddleware("admin", "manager")
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Lấy role từ context — được set bởi AuthMiddleware
		role, exists := c.Get("role")

		if !exists {
			responses.Error(c, http.StatusForbidden, "Role not found")
			c.Abort()
			return
		}

		userRole := role.(string)

		// Kiểm tra role có được phép không
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		responses.Error(c, http.StatusForbidden, "You do not have permission")
		c.Abort()
	}
}