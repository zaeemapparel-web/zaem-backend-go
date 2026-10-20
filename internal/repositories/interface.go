package repositories

import (
	"context"

	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
)

// ==================== USER ====================

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByResetToken(ctx context.Context, token string) (*models.User, error)
	FindByEmailVerifyToken(ctx context.Context, token string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

// ==================== CATEGORY ====================

type CategoryRepository interface {
	Create(ctx context.Context, cat *models.Category) error
	FindByID(ctx context.Context, id string) (*models.Category, error)
	FindBySlug(ctx context.Context, slug string) (*models.Category, error)
	FindAll(ctx context.Context) ([]models.Category, error)
	FindTree(ctx context.Context) ([]models.Category, error)
	Update(ctx context.Context, cat *models.Category) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	CountProducts(ctx context.Context, categoryID string) (int64, error)
	CountChildren(ctx context.Context, categoryID string) (int64, error)
}

// ==================== PRODUCT ====================

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product) error
	FindByID(ctx context.Context, id string) (*models.Product, error)
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	FindAll(ctx context.Context, filter dtos.ProductFilter, pagination dtos.Pagination) ([]models.Product, int64, error)
	FindFeatured(ctx context.Context, limit int) ([]models.Product, error)
	FindRelated(ctx context.Context, categoryID, excludeID string, limit int) ([]models.Product, error)
	Update(ctx context.Context, p *models.Product) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	IncrementViewCount(ctx context.Context, id string) error
}

// ==================== CART ====================

type CartRepository interface {
	FindOrCreateByUserID(ctx context.Context, userID string) (*models.Cart, error)
	Reload(ctx context.Context, userID string) (*models.Cart, error)
	FindItemByID(ctx context.Context, itemID string) (*models.CartItem, error)
	FindItemByProduct(ctx context.Context, cartID, productID string, size, color *string) (*models.CartItem, error)
	CreateItem(ctx context.Context, item *models.CartItem) error
	UpdateItem(ctx context.Context, item *models.CartItem) error
	UpdateItemQuantity(ctx context.Context, itemID string, qty int) error
	DeleteItem(ctx context.Context, itemID string) error
	ClearCart(ctx context.Context, cartID string) error
}

// ==================== WISHLIST ====================

type WishlistRepository interface {
	FindByUserID(ctx context.Context, userID string) ([]models.WishlistItem, error)
	FindByUserAndProduct(ctx context.Context, userID, productID string) (*models.WishlistItem, error)
	Create(ctx context.Context, item *models.WishlistItem) error
	Delete(ctx context.Context, userID, productID string) error
	DeleteByID(ctx context.Context, id string) error
	CountByUser(ctx context.Context, userID string) (int64, error)
}

// ==================== ADDRESS ====================

type AddressRepository interface {
	FindByUserID(ctx context.Context, userID string) ([]models.Address, error)
	FindByID(ctx context.Context, id string) (*models.Address, error)
	Create(ctx context.Context, addr *models.Address) error
	Update(ctx context.Context, addr *models.Address) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	ClearDefaultForUser(ctx context.Context, userID string) error
	CountByUser(ctx context.Context, userID string) (int64, error)
}

// ==================== ORDER ====================

type OrderRepository interface {
	CreateTx(ctx context.Context, order *models.Order) error
	FindByID(ctx context.Context, id string) (*models.Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	FindByUserID(ctx context.Context, userID string, pagination dtos.Pagination) ([]models.Order, int64, error)
	UpdateStatus(ctx context.Context, id string, fields map[string]interface{}) error
	GetAllAdmin(ctx context.Context, limit, offset int, status string) ([]models.Order, int64, error)
}  

// ==================== COUPON ====================

type CouponRepository interface {
	FindByCode(ctx context.Context, code string) (*models.Coupon, error)
	IncrementUsage(ctx context.Context, id string) error
	CreateUsage(ctx context.Context, usage *models.CouponUsage) error
	CountUsageByUser(ctx context.Context, couponID, userID string) (int64, error)
	CountUserOrders(ctx context.Context, userID string) (int64, error)
}
// ==================== ADMIN ====================

type AdminRepository interface {
	GetTotalRevenue(ctx context.Context) (float64, error)
	GetTotalOrders(ctx context.Context) (int64, error)
	GetTotalUsers(ctx context.Context) (int64, error)
	GetTotalProducts(ctx context.Context) (int64, error)
	GetPendingOrders(ctx context.Context) (int64, error)
	GetTodayRevenue(ctx context.Context) (float64, error)
	GetTodayOrders(ctx context.Context) (int64, error)
	GetLowStockCount(ctx context.Context, threshold int) (int64, error)
	GetAllOrders(ctx context.Context, limit, offset int, status string) ([]models.Order, int64, error)
	GetRecentOrders(ctx context.Context, limit int) ([]models.Order, error)
	GetTopProducts(ctx context.Context, limit int) ([]models.Product, error)
	
	GetAllUsers(ctx context.Context, limit, offset int, search string) ([]models.User, int64, error)
	FindUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUserFields(ctx context.Context, id string, fields map[string]interface{}) error
}
// ==================== REVIEW ====================

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	FindByID(ctx context.Context, id string) (*models.Review, error)
	FindByUserAndProduct(ctx context.Context, userID, productID string) (*models.Review, error)
	FindByProductID(ctx context.Context, productID string, page, limit int, onlyApproved bool) ([]models.Review, int64, error)
	FindAll(ctx context.Context, page, limit int, isApproved *bool) ([]models.Review, int64, error)
	Update(ctx context.Context, review *models.Review) error
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	FindHelpful(ctx context.Context, reviewID, userID string) (*models.ReviewHelpful, error)
	CreateHelpful(ctx context.Context, h *models.ReviewHelpful) error
	DeleteHelpful(ctx context.Context, id string) error
	IncrementHelpfulCount(ctx context.Context, reviewID string, delta int) error
	GetProductRatingStats(ctx context.Context, productID string) (float64, int64, error)
	HasPurchased(ctx context.Context, userID, productID string) (bool, error)
}

// ==================== NOTIFICATION ====================

type NotificationRepository interface {
	FindByUserID(ctx context.Context, userID string, page, limit int) ([]models.Notification, int64, error)
	CountUnread(ctx context.Context, userID string) (int64, error)
	MarkAsRead(ctx context.Context, id, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id, userID string) error
	Create(ctx context.Context, n *models.Notification) error
}

// ==================== NEWSLETTER ====================

type NewsletterRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.Newsletter, error)
	Create(ctx context.Context, n *models.Newsletter) error
	UpdateActive(ctx context.Context, id string, isActive bool) error
}