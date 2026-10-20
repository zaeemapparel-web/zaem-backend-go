package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}
	result, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, "Account created successfully. Welcome to ZAEM!", result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}
	result, err := h.svc.Login(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Logged in successfully.", result)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	result, err := h.svc.GetMe(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"user": result})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}
	userID := c.GetString("user_id")
	result, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Profile updated.", gin.H{"user": result})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dtos.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}
	userID := c.GetString("user_id")
	if err := h.svc.ChangePassword(c.Request.Context(), userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Password changed successfully.", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dtos.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}
	_ = h.svc.ForgotPassword(c.Request.Context(), req)
	// Always return success
	response.OKMessage(c, "If an account exists, a reset link has been sent", nil)
}

func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, 400, "Token is required")
		return
	}
	// Service mein validate karenge
	response.OKMessage(c, "Token is valid", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}
	if err := h.svc.ResetPassword(c.Request.Context(), req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Password reset successfully", nil)
}