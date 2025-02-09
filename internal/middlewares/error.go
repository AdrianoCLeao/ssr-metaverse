package middlewares

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "ssr-metaverse/internal/core/error"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() 

        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            log.Printf("Erro: %s", err.Error())

            c.JSON(http.StatusInternalServerError, error.Error(http.StatusInternalServerError, "Internal Server Error"))
            c.Abort()
        }
    }
}
