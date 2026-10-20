package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type ReviewHandler struct {
	svc *services.ReviewService
}

func NewReviewHandler(svc *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// ==================== USER / PUBLIC ====================

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID := c.Param("productId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, total, err := h.svc.GetProductReviews(c.Request.Context(), productID, page, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.ListPaginated(c, gin.H{"reviews": result}, page, limit, int(total))
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req dtos.CreateReviewRequest
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
	response.Created(c, "Review submitted for approval", gin.H{"review": result})
}

func (h *ReviewHandler) GetMyReviews(c *gin.Context) {
	userID := c.GetString("user_id")
	result, err := h.svc.GetMyReviews(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, len(result), gin.H{"reviews": result})
}

func (h *ReviewHandler) MarkHelpful(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.svc.MarkHelpful(c.Request.Context(), userID, id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Toggled helpful", nil)
}

// ==================== ADMIN ====================

func (h *ReviewHandler) GetAdminReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	var isApproved *bool
	if a := c.Query("isApproved"); a != "" {
		b := a == "true"
		isApproved = &b
	}

	result, total, err := h.svc.GetAllAdmin(c.Request.Context(), page, limit, isApproved)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.ListPaginated(c, gin.H{"reviews": result}, page, limit, int(total))
}

func (h *ReviewHandler) GetReviewStats(c *gin.Context) {
	stats, err := h.svc.GetReviewStats(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"stats": stats})
}

func (h *ReviewHandler) ApproveReview(c *gin.Context) {
	id := c.Param("id")
	isApproved := true
	req := dtos.AdminReviewActionRequest{IsApproved: &isApproved}

	if err := h.svc.AdminAction(c.Request.Context(), id, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Review approved", nil)
}

func (h *ReviewHandler) RejectReview(c *gin.Context) {
	id := c.Param("id")
	isApproved := false
	req := dtos.AdminReviewActionRequest{IsApproved: &isApproved}

	if err := h.svc.AdminAction(c.Request.Context(), id, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Review rejected", nil)
}

func (h *ReviewHandler) AdminReply(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Reply string `json:"reply" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, 400, "Reply is required")
		return
	}

	req := dtos.AdminReviewActionRequest{AdminReply: &body.Reply}
	if err := h.svc.AdminAction(c.Request.Context(), id, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Reply added", nil)
}

func (h *ReviewHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	isAdmin := c.GetString("user_role") == "ADMIN"

	if err := h.svc.Delete(c.Request.Context(), userID, isAdmin, id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Review deleted", nil)
}