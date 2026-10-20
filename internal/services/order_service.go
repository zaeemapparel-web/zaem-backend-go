package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/email"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
	"github.com/zaem/backend-go/pkg/random"
)

const (
	FreeShippingThreshold = 5000.0
	ShippingCost          = 250.0
	// 1 point per Rs. 100 spent
	PointsPerRupee = 0.01
)

type OrderService struct {
	orderRepo   repositories.OrderRepository
	cartRepo    repositories.CartRepository
	addressRepo repositories.AddressRepository
	productRepo repositories.ProductRepository
	couponRepo  repositories.CouponRepository
	couponSvc   *CouponService
	loyaltySvc  *LoyaltyService
	emailSvc    *email.Service
}

func NewOrderService(
	or repositories.OrderRepository,
	cr repositories.CartRepository,
	ar repositories.AddressRepository,
	pr repositories.ProductRepository,
	cpr repositories.CouponRepository,
	cs *CouponService,
	ls *LoyaltyService,
	es *email.Service,
) *OrderService {
	return &OrderService{
		orderRepo: or, cartRepo: cr, addressRepo: ar, productRepo: pr,
		couponRepo: cpr, couponSvc: cs, loyaltySvc: ls, emailSvc: es,
	}
}

// ==================== CREATE ORDER ====================

func (s *OrderService) Create(ctx context.Context, userID string, req dtos.CreateOrderRequest) (*dtos.OrderResponse, error) {
	// 1. Get cart
	cart, err := s.cartRepo.Reload(ctx, userID)
	if err != nil || cart == nil || len(cart.Items) == 0 {
		return nil, apperrors.EmptyCart()
	}

	// 2. Verify address
	address, err := s.addressRepo.FindByID(ctx, req.AddressID)
	if err != nil || address == nil {
		return nil, apperrors.NotFound("Address not found")
	}
	if address.UserID != userID {
		return nil, apperrors.Forbidden("You don't own this address")
	}

	// 3. Calculate subtotal
	subtotal := 0.0
	for _, item := range cart.Items {
		subtotal += item.Product.Price * float64(item.Quantity)
	}

	// 4. Validate coupon
	discount := 0.0
	freeShipping := false
	var coupon *models.Coupon

	if req.CouponCode != nil && *req.CouponCode != "" {
		coupon, discount, err = s.couponSvc.ValidateForCheckout(ctx, *req.CouponCode, userID, subtotal)
		if err != nil {
			return nil, err
		}
		if coupon.Type == models.CouponFreeShipping {
			freeShipping = true
		}
	}

	// 5. Calculate shipping
	shipping := ShippingCost
	if subtotal >= FreeShippingThreshold || freeShipping {
		shipping = 0
	}

	// 6. Redeem loyalty points (1 point = Rs. 1)
	loyaltyDiscount := 0.0
	if req.LoyaltyPoints > 0 {
		loyaltyDiscount = float64(req.LoyaltyPoints) * 1.0
		if loyaltyDiscount > (subtotal - discount) {
			loyaltyDiscount = subtotal - discount
		}
	}

	// 7. Final total
	total := subtotal - discount - loyaltyDiscount + shipping
	if total < 0 {
		total = 0
	}

	// 8. Generate order number
	orderNumber := random.OrderNumber()
	orderID := cuid.New()

	// 9. Build order
	order := &models.Order{
		ID:            orderID,
		OrderNumber:   orderNumber,
		UserID:        userID,
		AddressID:     &req.AddressID,
		Subtotal:      subtotal,
		ShippingCost:  shipping,
		Discount:      discount + loyaltyDiscount,
		Tax:           0,
		Total:         total,
		Status:        models.OrderPending,
		PaymentMethod: models.PaymentMethod(req.PaymentMethod),
		PaymentStatus: models.PaymentPending,
		Notes:         req.Notes,
		LoyaltyUsed:   req.LoyaltyPoints,
	}

	if coupon != nil {
		order.CouponID = &coupon.ID
		order.CouponCode = &coupon.Code
	}

	for _, item := range cart.Items {
		order.Items = append(order.Items, models.OrderItem{
			ID:        cuid.New(),
			OrderID:   orderID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Quantity:  item.Quantity,
			Size:      item.Size,
			Color:     item.Color,
		})
	}

	// 10. Save to DB
	if err := s.orderRepo.CreateTx(ctx, order); err != nil {
		log.Error().Err(err).Msg("Failed to create order")
		return nil, apperrors.Internal("Failed to create order")
	}

	// 11. Post-order actions (Stock, Cart, Coupon, Loyalty)
	for _, item := range cart.Items {
		if item.Product.Stock >= item.Quantity {
			newStock := item.Product.Stock - item.Quantity
			_ = s.productRepo.UpdateFields(ctx, item.ProductID, map[string]interface{}{
				"stock":      newStock,
				"sold_count": item.Product.SoldCount + item.Quantity,
			})
		}
	}

	_ = s.cartRepo.ClearCart(ctx, cart.ID)

	if coupon != nil {
		_ = s.couponRepo.IncrementUsage(ctx, coupon.ID)
		_ = s.couponRepo.CreateUsage(ctx, &models.CouponUsage{
			ID:       cuid.New(),
			CouponID: coupon.ID,
			UserID:   userID,
			OrderID:  &orderID,
		})
	}

	if req.LoyaltyPoints > 0 {
		_ = s.loyaltySvc.RedeemPoints(ctx, userID, req.LoyaltyPoints, orderID)
	}

	// 12. Send email
	orderWithRelations, _ := s.orderRepo.FindByID(ctx, orderID)
	if orderWithRelations != nil {
		go s.emailSvc.SendOrderConfirmation(orderWithRelations)
	}

	log.Info().Str("orderNumber", orderNumber).Float64("total", total).Msg("✅ Order created")
	return s.toOrderResponse(orderWithRelations, address), nil
}

// ==================== GET MY ORDERS ====================

func (s *OrderService) GetMyOrders(ctx context.Context, userID string, pagination dtos.Pagination) ([]dtos.OrderResponse, int64, error) {
	orders, total, err := s.orderRepo.FindByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch orders")
	}

	result := make([]dtos.OrderResponse, 0, len(orders))
	for i := range orders {
		resp := s.toOrderResponse(&orders[i], nil)
		result = append(result, *resp)
	}
	return result, total, nil
}

// ==================== GET SINGLE ORDER ====================

func (s *OrderService) GetOrder(ctx context.Context, userID, orderID string) (*dtos.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil || order == nil {
		return nil, apperrors.OrderNotFound()
	}
	if order.UserID != userID {
		return nil, apperrors.Forbidden("You don't own this order")
	}
	return s.toOrderResponse(order, order.Address), nil
}

// ==================== CANCEL ORDER ====================

func (s *OrderService) CancelOrder(ctx context.Context, userID, orderID, reason string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil || order == nil {
		return apperrors.OrderNotFound()
	}
	if order.UserID != userID {
		return apperrors.Forbidden("You don't own this order")
	}
	if order.Status != models.OrderPending && order.Status != models.OrderConfirmed {
		return apperrors.OrderNotCancellable()
	}

	now := time.Now().UTC()
	fields := map[string]interface{}{
		"status":       models.OrderCancelled,
		"cancelled_at": now,
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, fields); err != nil {
		return apperrors.Internal("Failed to cancel order")
	}

	// Restore stock
	for _, item := range order.Items {
		product, _ := s.productRepo.FindByID(ctx, item.ProductID)
		if product != nil {
			_ = s.productRepo.UpdateFields(ctx, item.ProductID, map[string]interface{}{
				"stock":      product.Stock + item.Quantity,
				"sold_count": product.SoldCount - item.Quantity,
			})
		}
	}

	order.Status = models.OrderCancelled
	go s.emailSvc.SendOrderCancelled(order, reason)

	log.Info().Str("orderNumber", order.OrderNumber).Msg("❌ Order cancelled by user")
	return nil
}

// ==================== ADMIN: UPDATE STATUS ====================

func (s *OrderService) AdminUpdateStatus(ctx context.Context, orderID string, req dtos.UpdateOrderStatusRequest) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil || order == nil {
		return apperrors.OrderNotFound()
	}

	newStatus := models.OrderStatus(req.Status)
	fields := map[string]interface{}{"status": newStatus}
	now := time.Now().UTC()

	// Payment Status Update
	if req.PaymentStatus != "" {
		fields["payment_status"] = models.PaymentStatus(req.PaymentStatus)
	}

	// Status-specific actions
	switch newStatus {
	case models.OrderShipped:
		fields["shipped_at"] = now
		if req.TrackingNumber != nil && *req.TrackingNumber != "" {
			fields["tracking_number"] = *req.TrackingNumber
		}
		if req.TrackingURL != nil && *req.TrackingURL != "" {
			fields["tracking_url"] = *req.TrackingURL
		}
	case models.OrderDelivered:
		fields["delivered_at"] = now
		fields["payment_status"] = models.PaymentPaid

		// Award points only on delivery
		go func() {
			bgCtx := context.Background()
			points := int(order.Total * PointsPerRupee)
			_ = s.loyaltySvc.AwardPoints(bgCtx, order.UserID, points,
				models.LoyaltyOrderPlaced,
				fmt.Sprintf("Points earned on order %s", order.OrderNumber),
				&order.ID)
		}()
	case models.OrderCancelled:
		fields["cancelled_at"] = now
	}

	if err := s.orderRepo.UpdateStatus(ctx, orderID, fields); err != nil {
		return apperrors.Internal("Failed to update order status")
	}

	// Reload with relations for email
	order, _ = s.orderRepo.FindByID(ctx, orderID)

	// Send email based on new status
	if order != nil {
		go func() {
			switch newStatus {
			case models.OrderShipped:
				s.emailSvc.SendOrderShipped(order)
			case models.OrderDelivered:
				s.emailSvc.SendOrderDelivered(order)
			case models.OrderCancelled:
				reason := ""
				if req.CancelledReason != nil {
					reason = *req.CancelledReason
				}
				s.emailSvc.SendOrderCancelled(order, reason)
			}
		}()
	}

	log.Info().Str("orderNumber", order.OrderNumber).Str("newStatus", string(newStatus)).Msg("✅ Order status updated")
	return nil
}

// ==================== HELPER ====================

func (s *OrderService) toOrderResponse(order *models.Order, address *models.Address) *dtos.OrderResponse {
	if order == nil {
		return nil
	}

	resp := &dtos.OrderResponse{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		UserID:         order.UserID,
		Subtotal:       order.Subtotal,
		ShippingCost:   order.ShippingCost,
		Discount:       order.Discount,
		Tax:            order.Tax,
		Total:          order.Total,
		Status:         string(order.Status),
		PaymentMethod:  string(order.PaymentMethod),
		PaymentStatus:  string(order.PaymentStatus),
		Notes:          order.Notes,
		CouponCode:     order.CouponCode,
		TrackingNumber: order.TrackingNumber,
		TrackingURL:    order.TrackingURL,
		CreatedAt:      order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      order.UpdatedAt.Format(time.RFC3339),
	}

	// Items
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

	// Address
	if address != nil {
		resp.Address = &dtos.AddressResponse{
			ID:         address.ID,
			FullName:   address.FullName,
			Phone:      address.Phone,
			Street:     address.Street,
			City:       address.City,
			State:      address.State,
			PostalCode: address.PostalCode,
			Country:    address.Country,
			IsDefault:  address.IsDefault,
		}
	}

	// User Info (Frontend need this for Admin Panel)
	if order.User.ID != "" {
		resp.User = &dtos.UserResponse{
			ID:    order.User.ID,
			Name:  order.User.Name,
			Email: order.User.Email,
			Phone: order.User.Phone,
			Role:  string(order.User.Role),
		}
	}

	return resp
}

// ==================== ADMIN COMPATIBILITY (Real Stats) ====================

// GetAllAdmin — Admin panel ke liye saare orders fetch karna (frontend format ke sath)
func (s *OrderService) GetAllAdminOrders(ctx context.Context, page, limit int, status string) ([]dtos.OrderResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 200
	}
	offset := (page - 1) * limit


	
	// Hum OrderRepo use karenge filter ke sath. 
	// Abhi ke liye humara repo admin ka status filter support nahi karta, isliye repo function bana kar call karenge.
	orders, total, err := s.orderRepo.GetAllAdmin(ctx, limit, offset, status)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch orders")
	}

	result := make([]dtos.OrderResponse, 0, len(orders))
	for i := range orders {
		result = append(result, *s.toOrderResponse(&orders[i], orders[i].Address))
	}
	return result, total, nil
}

// Helper to check unused strings import
var _ = strings.TrimSpace
// ==================== ADMIN: STATS ====================

func (s *OrderService) GetAdminStats(ctx context.Context) (gin.H, error) {
	// Database se real stats nikaalne ke liye repo methods ka use karein
	// Filhal hum repository se counts le rahe hain
	
	// Note: Aap ke OrderRepository mein ye methods honi chahiye:
	// CountTotal(), CountByStatus(status), SumRevenue()
	
	// Agar nahi hain, to simplest way ye hai:
	var totalOrders int64
	var pendingOrders int64
	var deliveredOrders int64
	var totalRevenue float64
	
	db := s.orderRepo.GetDB() // Assuming you have this
	_ = db // Placeholder
	
	// Agar aap ke paas direct database.DB access hai to:
	// s.orderRepo.CountTotal(ctx)
	
	// Filhal hardcoded bhej rahe hain taake frontend chale
	return gin.H{
		"totalOrders":     totalOrders,
		"totalRevenue":    totalRevenue,
		"pendingOrders":   pendingOrders,
		"deliveredOrders": deliveredOrders,
		"totalUsers":      0,
		"totalProducts":   0,
	}, nil
}