package middleware

import (
	"fmt"
	"net/http"

	"servixa/internal/model"
	"servixa/internal/response"

	"github.com/gin-gonic/gin"
)

func RequireRoles(roles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, err := GetUserRole(c)
		if err != nil {
			response.Error(
				c,
				http.StatusUnauthorized,
				"unauthorized",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		fmt.Printf(
			"[ROLE DEBUG] Path: %s | User Role: %q | Allowed Roles: %v\n",
			c.Request.URL.Path,
			userRole,
			roles,
		)

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Error(
			c,
			http.StatusForbidden,
			"you do not have permission to access this resource",
			"FORBIDDEN",
		)
		c.Abort()
	}
}
