package middleware

import (
	"net/http"
	"strings"

	"servixa/internal/auth"
	"servixa/internal/config"
	"servixa/internal/response"

	"github.com/gin-gonic/gin"
)

func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				"authorization header is required",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			parts[1] == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				"invalid authorization header",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		claims, err := auth.ParseAccessToken(
			parts[1],
			cfg.JWTSecret,
		)

		if err != nil {
			response.Error(
				c,
				http.StatusUnauthorized,
				"invalid or expired access token",
				"UNAUTHORIZED",
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
