package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerCartRoutes(r *gin.RouterGroup, h *handlers.CartHandler) {
	cart := r.Group("/cart")
	cart.Use(middleware.Protect())

	cart.GET("/", h.Get)
	cart.POST("/", h.Add)
	cart.PUT("/:itemId", h.UpdateQuantity)
	cart.DELETE("/:itemId", h.Remove)
	cart.DELETE("/", h.Clear)
}