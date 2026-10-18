package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type NewsletterHandler struct {
	svc *services.NewsletterService
}

func NewNewsletterHandler(svc *services.NewsletterService) *NewsletterHandler {
	return &NewsletterHandler{svc: svc}
}

func (h *NewsletterHandler) Subscribe(c *gin.Context) {
	var req dtos.NewsletterSubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	if err := h.svc.Subscribe(c.Request.Context(), req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Subscribed successfully", nil)
}

func (h *NewsletterHandler) Unsubscribe(c *gin.Context) {
	var req dtos.NewsletterUnsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	if err := h.svc.Unsubscribe(c.Request.Context(), req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Unsubscribed successfully", nil)
}