package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RequestLogger logs every HTTP request with structured fields
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger := log.Info()
		if status >= 500 {
			logger = log.Error()
		} else if status >= 400 {
			logger = log.Warn()
		}

		requestID, _ := c.Get(RequestIDKey)

		logger.
			Str("requestId", toString(requestID)).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("userAgent", c.Request.UserAgent()).
			Msg("📨 HTTP Request")
	}
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}