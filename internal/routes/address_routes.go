package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerAddressRoutes(r *gin.RouterGroup, h *handlers.AddressHandler) {
	addresses := r.Group("/addresses")
	addresses.Use(middleware.Protect())

	addresses.GET("", h.GetMy)
	addresses.POST("", h.Create)
	addresses.PUT("/:id", h.Update)
	addresses.DELETE("/:id", h.Delete)
	addresses.PATCH("/:id/default", h.SetDefault)
}