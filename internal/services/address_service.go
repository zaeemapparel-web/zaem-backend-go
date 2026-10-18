package services

import (
	"context"
	"strings"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
)

type AddressService struct {
	repo repositories.AddressRepository
}

func NewAddressService(r repositories.AddressRepository) *AddressService {
	return &AddressService{repo: r}
}

func (s *AddressService) GetMyAddresses(ctx context.Context, userID string) ([]dtos.AddressResponse, error) {
	addresses, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch addresses")
	}

	resp := make([]dtos.AddressResponse, 0, len(addresses))
	for i := range addresses {
		resp = append(resp, s.toResponse(&addresses[i]))
	}
	return resp, nil
}

func (s *AddressService) Create(ctx context.Context, userID string, req dtos.CreateAddressRequest) (*dtos.AddressResponse, error) {
	// Check address count (max 10)
	count, _ := s.repo.CountByUser(ctx, userID)
	if count >= 10 {
		return nil, apperrors.BadRequest("Maximum 10 addresses allowed")
	}

	// If this is first address, make it default
	isDefault := req.IsDefault
	if count == 0 {
		isDefault = true
	}

	// If setting as default, clear others
	if isDefault {
		_ = s.repo.ClearDefaultForUser(ctx, userID)
	}

	country := strings.TrimSpace(req.Country)
	if country == "" {
		country = "Pakistan"
	}

	addr := &models.Address{
		ID:         cuid.New(),
		UserID:     userID,
		FullName:   strings.TrimSpace(req.FullName),
		Phone:      strings.TrimSpace(req.Phone),
		Street:     strings.TrimSpace(req.Street),
		City:       strings.TrimSpace(req.City),
		State:      req.State,
		PostalCode: strings.TrimSpace(req.PostalCode),
		Country:    country,
		IsDefault:  isDefault,
	}

	if err := s.repo.Create(ctx, addr); err != nil {
		return nil, apperrors.Internal("Failed to create address")
	}

	resp := s.toResponse(addr)
	return &resp, nil
}

func (s *AddressService) Update(ctx context.Context, userID, id string, req dtos.UpdateAddressRequest) (*dtos.AddressResponse, error) {
	addr, err := s.repo.FindByID(ctx, id)
	if err != nil || addr == nil {
		return nil, apperrors.NotFound("Address not found")
	}
	if addr.UserID != userID {
		return nil, apperrors.Forbidden("You don't own this address")
	}

	updates := map[string]interface{}{}

	if req.FullName != nil {
		updates["full_name"] = strings.TrimSpace(*req.FullName)
	}
	if req.Phone != nil {
		updates["phone"] = strings.TrimSpace(*req.Phone)
	}
	if req.Street != nil {
		updates["street"] = strings.TrimSpace(*req.Street)
	}
	if req.City != nil {
		updates["city"] = strings.TrimSpace(*req.City)
	}
	if req.State != nil {
		updates["state"] = *req.State
	}
	if req.PostalCode != nil {
		updates["postal_code"] = strings.TrimSpace(*req.PostalCode)
	}
	if req.Country != nil {
		updates["country"] = strings.TrimSpace(*req.Country)
	}

	if req.IsDefault != nil {
		if *req.IsDefault {
			_ = s.repo.ClearDefaultForUser(ctx, userID)
			updates["is_default"] = true
		} else {
			updates["is_default"] = false
		}
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateFields(ctx, id, updates); err != nil {
			return nil, apperrors.Internal("Failed to update address")
		}
	}

	addr, _ = s.repo.FindByID(ctx, id)
	resp := s.toResponse(addr)
	return &resp, nil
}

func (s *AddressService) Delete(ctx context.Context, userID, id string) error {
	addr, err := s.repo.FindByID(ctx, id)
	if err != nil || addr == nil {
		return apperrors.NotFound("Address not found")
	}
	if addr.UserID != userID {
		return apperrors.Forbidden("You don't own this address")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.Internal("Failed to delete address")
	}
	return nil
}

func (s *AddressService) SetDefault(ctx context.Context, userID, id string) error {
	addr, err := s.repo.FindByID(ctx, id)
	if err != nil || addr == nil {
		return apperrors.NotFound("Address not found")
	}
	if addr.UserID != userID {
		return apperrors.Forbidden("You don't own this address")
	}

	if err := s.repo.ClearDefaultForUser(ctx, userID); err != nil {
		return apperrors.Internal("Failed to set default")
	}

	if err := s.repo.UpdateFields(ctx, id, map[string]interface{}{"is_default": true}); err != nil {
		return apperrors.Internal("Failed to set default")
	}

	return nil
}

func (s *AddressService) toResponse(a *models.Address) dtos.AddressResponse {
	return dtos.AddressResponse{
		ID:         a.ID,
		FullName:   a.FullName,
		Phone:      a.Phone,
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		PostalCode: a.PostalCode,
		Country:    a.Country,
		IsDefault:  a.IsDefault,
		CreatedAt:  a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  a.UpdatedAt.Format(time.RFC3339),
	}
}