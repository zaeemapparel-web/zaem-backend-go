package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
)

type NotificationHandler struct {
	svc *services.NotificationService
}

func NewNotificationHandler(svc *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) GetMy(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, total, unread, err := h.svc.GetMyNotifications(c.Request.Context(), userID, page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"notifications": result},
		"unread":  unread,
		"pagination": gin.H{
			"page": page, "limit": limit, "total": total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.svc.MarkAsRead(c.Request.Context(), userID, id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Marked as read", nil)
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetString("user_id")
	if err := h.svc.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "All marked as read", nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.svc.Delete(c.Request.Context(), userID, id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Notification deleted", nil)
}