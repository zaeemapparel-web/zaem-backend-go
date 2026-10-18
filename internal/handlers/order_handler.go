package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type OrderHandler struct {
	svc *services.OrderService
}

func NewOrderHandler(svc *services.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// Create — user places an order
func (h *OrderHandler) Create(c *gin.Context) {
	var req dtos.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID := c.GetString("user_id")
	result, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Order placed successfully", gin.H{"order": result})
}

// GetMyOrders — user's order history
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetString("user_id")

	pagination := dtos.NewPagination(
		c.Query("page"),
		c.Query("limit"),
		"created_at",
		"desc",
	)

	orders, total, err := h.svc.GetMyOrders(c.Request.Context(), userID, pagination)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.ListPaginated(c, gin.H{"orders": orders}, pagination.Page, pagination.Limit, int(total))
}

// GetByID — single order details
func (h *OrderHandler) GetByID(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")

	result, err := h.svc.GetOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, gin.H{"order": result})
}

// Cancel — user cancels their own order
func (h *OrderHandler) Cancel(c *gin.Context) {
	orderID := c.Param("id")

	var req dtos.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID := c.GetString("user_id")
	if err := h.svc.CancelOrder(c.Request.Context(), userID, orderID, req.Reason); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKMessage(c, "Order cancelled successfully", nil)
}

// AdminUpdateStatus — admin changes status (triggers email)
func (h *OrderHandler) AdminUpdateStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req dtos.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	if err := h.svc.AdminUpdateStatus(c.Request.Context(), orderID, req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKMessage(c, "Order status updated successfully", nil)
}