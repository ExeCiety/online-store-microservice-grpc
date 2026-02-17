package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func OK(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{Success: true, Message: message, Data: data})
}

func Fail(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, APIResponse{Success: false, Message: message, Error: err})
}
