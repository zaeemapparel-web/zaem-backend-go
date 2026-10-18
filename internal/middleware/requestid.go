package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/pkg/cuid"
)

const RequestIDKey = "request_id"

// RequestID adds a unique ID to every request (for tracing)
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = cuid.New()
		}

		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}