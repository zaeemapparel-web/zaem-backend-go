package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type WishlistHandler struct {
	svc *services.WishlistService
}

func NewWishlistHandler(svc *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{svc: svc}
}

func (h *WishlistHandler) GetMy(c *gin.Context) {
	userID := c.GetString("user_id")
	result, err := h.svc.GetMyWishlist(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, len(result), gin.H{"items": result})
}

func (h *WishlistHandler) Add(c *gin.Context) {
	var req dtos.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID := c.GetString("user_id")
	if err := h.svc.Add(c.Request.Context(), userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Added to wishlist", nil)
}

func (h *WishlistHandler) Remove(c *gin.Context) {
	productID := c.Param("productId")
	userID := c.GetString("user_id")

	if err := h.svc.Remove(c.Request.Context(), userID, productID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Removed from wishlist", nil)
}