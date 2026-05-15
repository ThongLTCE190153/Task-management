package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		fmt.Println("Method:", c.Request.Method)
		fmt.Println("Path:", c.Request.URL.Path)
		fmt.Println("Status:", c.Writer.Status())
		fmt.Println("Duration:", duration)
		fmt.Println("----------------------")
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}