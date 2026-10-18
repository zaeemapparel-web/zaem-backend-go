package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type CartHandler struct {
	svc *services.CartService
}

func NewCartHandler(svc *services.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	result, err := h.svc.GetCart(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"cart": result})
}

func (h *CartHandler) Add(c *gin.Context) {
	var req dtos.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID := c.GetString("user_id")
	result, err := h.svc.AddItem(c.Request.Context(), userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Item added to cart", gin.H{"cart": result})
}

func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	itemID := c.Param("itemId")

	var req dtos.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID := c.GetString("user_id")
	result, err := h.svc.UpdateQuantity(c.Request.Context(), userID, itemID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Cart updated", gin.H{"cart": result})
}

func (h *CartHandler) Remove(c *gin.Context) {
	itemID := c.Param("itemId")
	userID := c.GetString("user_id")

	result, err := h.svc.RemoveItem(c.Request.Context(), userID, itemID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Item removed from cart", gin.H{"cart": result})
}

func (h *CartHandler) Clear(c *gin.Context) {
	userID := c.GetString("user_id")
	result, err := h.svc.ClearCart(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Cart cleared", gin.H{"cart": result})
}