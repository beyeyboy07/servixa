package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"

	"servixa/internal/model"
)

func GetUserID(c *gin.Context) (uint, error) {
	value, exists := c.Get("user_id")

	if !exists {
		return 0, errors.New("user id not found in context")
	}

	userID, ok := value.(uint)

	if !ok {
		return 0, errors.New("invalid user id in context")
	}

	return userID, nil
}

func GetUserRole(c *gin.Context) (model.UserRole, error) {
	value, exists := c.Get("user_role")

	if !exists {
		return "", errors.New("user role not found in context")
	}

	role, ok := value.(model.UserRole)

	if !ok {
		return "", errors.New("invalid user role in context")
	}

	return role, nil
}
