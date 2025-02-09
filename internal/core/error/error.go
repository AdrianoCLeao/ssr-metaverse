package error

import (
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