package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerWishlistRoutes(r *gin.RouterGroup, h *handlers.WishlistHandler) {
	wishlist := r.Group("/wishlist")
	wishlist.Use(middleware.Protect())

	wishlist.GET("", h.GetMy)
	wishlist.POST("", h.Add)
	wishlist.DELETE("/:productId", h.Remove)
}