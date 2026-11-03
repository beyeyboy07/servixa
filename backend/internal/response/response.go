package response

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, code string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Message: message,
		Error:   code,
	})
}
