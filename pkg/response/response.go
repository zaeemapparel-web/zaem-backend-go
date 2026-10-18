package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/zaem/backend-go/internal/errors"
)

// ==================== SUCCESS RESPONSES ====================

// Success sends a standard success response
// Shape: { success: true, message: "...", data: {...} }
func Success(c *gin.Context, status int, message string, data interface{}) {
	body := gin.H{
		"success": true,
	}

	if message != "" {
		body["message"] = message
	}
	if data != nil {
		body["data"] = data
	}

	c.JSON(status, body)
}

// OK sends a 200 with data
func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, "", data)
}

// OKMessage sends a 200 with message + data
func OKMessage(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// Created sends 201 with data
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// ==================== LIST RESPONSES ====================

// List sends a list response with count
// Shape: { success: true, count: N, data: {...} }
func List(c *gin.Context, count int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   count,
		"data":    data,
	})
}

// ListPaginated sends paginated list
// Shape: { success: true, data: {...}, pagination: {...} }
func ListPaginated(c *gin.Context, data interface{}, page, limit, total int) {
	totalPages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// ==================== ERROR RESPONSES ====================

// Error sends a standard error response
// Shape: { success: false, message: "...", code?: "..." }
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

// ErrorWithCode sends error with code
func ErrorWithCode(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"code":    code,
		"message": message,
	})
}

// HandleError handles our AppError type or falls back to 500
// This is the MAIN error handler used in handlers
func HandleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError

	// Check if it's our custom error type
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, gin.H{
			"success": false,
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	// Unknown error → 500
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Internal server error",
	})
}