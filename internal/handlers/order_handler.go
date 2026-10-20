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
	response.Created(c, "Order placed successfully.", gin.H{"order": result})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	pagination := dtos.NewPagination(c.Query("page"), c.Query("limit"), "created_at", "desc")
	
	// GetMyOrders returns []OrderResponse, int64, error
	orders, total, err := h.svc.GetMyOrders(c.Request.Context(), userID, pagination)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, int(total), gin.H{"orders": orders})
}

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

func (h *OrderHandler) Cancel(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")
	if err := h.svc.CancelOrder(c.Request.Context(), userID, orderID, "User requested"); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Order cancelled.", nil)
}

// ==================== ADMIN ====================

func (h *OrderHandler) AdminGetAllOrders(c *gin.Context) {
	pagination := dtos.NewPagination(c.Query("page"), c.Query("limit"), "created_at", "desc")
	status := c.Query("status")
	
	orders, total, err := h.svc.GetAllAdminOrders(c.Request.Context(), pagination.Page, pagination.Limit, status)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"count":   len(orders),
		"pagination": gin.H{
			"total":      total,
			"page":       pagination.Page,
			"limit":      pagination.Limit,
			"totalPages": (int(total) + pagination.Limit - 1) / pagination.Limit,
		},
		"data": gin.H{"orders": orders},
	})
}

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
	response.OKMessage(c, "Order updated.", nil)
}

func (h *OrderHandler) AdminGetStats(c *gin.Context) {
	stats, err := h.svc.GetAdminStats(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"stats": stats})
}