package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// ==================== REGISTER ====================

func (h *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Account created successfully", result)
}

// ==================== LOGIN ====================

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OK(c, result)
}

// ==================== GET ME ====================

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	result, err := h.authService.GetMe(c.Request.Context(), userID.(string))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"user": result})
}

// ==================== UPDATE PROFILE ====================

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID, _ := c.Get("user_id")
	result, err := h.authService.UpdateProfile(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKMessage(c, "Profile updated successfully", gin.H{"user": result})
}

// ==================== CHANGE PASSWORD ====================

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dtos.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.authService.ChangePassword(c.Request.Context(), userID.(string), req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKMessage(c, "Password changed successfully", nil)
}

// ==================== FORGOT PASSWORD ====================

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dtos.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.OKMessage(c, "If an account exists with this email, a password reset link has been sent.", nil)
}

// ==================== RESET PASSWORD ====================

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req); err != nil {
		response.HandleError(c, err)
		return
	}

	log.Info().Msg("✅ Password reset endpoint completed")
	response.OKMessage(c, "Password reset successfully. You can now log in.", nil)
}