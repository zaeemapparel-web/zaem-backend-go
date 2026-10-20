package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type AdminHandler struct {
	svc *services.AdminService
}

func NewAdminHandler(svc *services.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	result, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"stats": result})
}

func (h *AdminHandler) GetRecentOrders(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.svc.GetRecentOrders(c.Request.Context(), limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"orders": result})
}

func (h *AdminHandler) GetTopProducts(c *gin.Context) {
	limit := 5
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.svc.GetTopProducts(c.Request.Context(), limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"products": result})
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	search := c.Query("search")

	result, total, err := h.svc.GetAllUsers(c.Request.Context(), page, limit, search)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.ListPaginated(c, gin.H{"users": result}, page, limit, int(total))
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req dtos.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	adminID := c.GetString("user_id")
	if err := h.svc.UpdateUserRole(c.Request.Context(), adminID, userID, req.Role); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKMessage(c, "User role updated", nil)
}

func (h *AdminHandler) ToggleUserActive(c *gin.Context) {
	userID := c.Param("id")

	var req dtos.ToggleUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	adminID := c.GetString("user_id")
	if err := h.svc.ToggleUserActive(c.Request.Context(), adminID, userID, req.IsActive); err != nil {
		response.HandleError(c, err)
		return
	}

	msg := "User activated"
	if !req.IsActive {
		msg = "User deactivated"
	}
	response.OKMessage(c, msg, nil)
}
func (s *AdminService) GetAllOrders(ctx context.Context, page, limit int, status string) ([]dtos.OrderResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	orders, total, err := s.repo.GetAllOrders(ctx, limit, offset, status)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch orders")
	}

	// OrderService ka toOrderResponse function use karne ke liye
	// kyunki AdminService ke paas wo helper nahi hai, hum yahan manual mapping karenge
	result := make([]dtos.OrderResponse, 0, len(orders))
	for i := range orders {
		order := &orders[i]
		resp := dtos.OrderResponse{
			ID:            order.ID,
			OrderNumber:   order.OrderNumber,
			UserID:        order.UserID,
			Subtotal:      order.Subtotal,
			ShippingCost:  order.ShippingCost,
			Discount:      order.Discount,
			Tax:           order.Tax,
			Total:         order.Total,
			Status:        string(order.Status),
			PaymentMethod: string(order.PaymentMethod),
			PaymentStatus: string(order.PaymentStatus),
			Notes:         order.Notes,
			CreatedAt:     order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     order.UpdatedAt.Format(time.RFC3339),
		}

		for _, item := range order.Items {
			resp.Items = append(resp.Items, dtos.OrderItemResponse{
				ID:        item.ID,
				ProductID: item.ProductID,
				Name:      item.Name,
				Price:     item.Price,
				Quantity:  item.Quantity,
				Size:      item.Size,
				Color:     item.Color,
				Image:     item.Image,
			})
		}
		
		result = append(result, resp)
	}
	return result, total, nil
}