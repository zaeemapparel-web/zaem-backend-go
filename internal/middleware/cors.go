package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/config"
)

// CORS allows requests from frontend origins
func CORS() gin.HandlerFunc {
	allowed := map[string]bool{
		"http://localhost:3000":            true,
		"http://localhost:3001":            true,
		"https://zaemstore.com":            true,
		"https://www.zaemstore.com":        true,
		"https://zaem-frontend.vercel.app": true,
	}
	if config.C.FrontendURL != "" {
		allowed[config.C.FrontendURL] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}