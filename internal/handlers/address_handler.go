package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type AddressHandler struct {
	svc *services.AddressService
}

func NewAddressHandler(svc *services.AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

func (h *AddressHandler) GetMy(c *gin.Context) {
	userID := c.GetString("user_id")
	result, err := h.svc.GetMyAddresses(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, len(result), gin.H{"addresses": result})
}

func (h *AddressHandler) Create(c *gin.Context) {
	var req dtos.CreateAddressRequest
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
	response.Created(c, "Address added", gin.H{"address": result})
}

func (h *AddressHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dtos.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	userID := c.GetString("user_id")
	result, err := h.svc.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Address updated", gin.H{"address": result})
}

func (h *AddressHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.svc.Delete(c.Request.Context(), userID, id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Address deleted", nil)
}

func (h *AddressHandler) SetDefault(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.svc.SetDefault(c.Request.Context(), userID, id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Default address set", nil)
}