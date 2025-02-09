package error

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func Error(code int, message string) APIError {
    return APIError{
        Code:    code,
        Message: message,
    }
}

func RespondWithError(c *gin.Context, err APIError) {
    c.JSON(err.Code, gin.H{"error": err.Message})
    c.Abort()
}

func (e *APIError) Error() string {
    return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}