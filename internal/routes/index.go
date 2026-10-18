package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/email"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/internal/services"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	// ==================== REPOSITORIES ====================
	userRepo := repositories.NewUserRepository()
	categoryRepo := repositories.NewCategoryRepository()
	productRepo := repositories.NewProductRepository()
	cartRepo := repositories.NewCartRepository()
	wishlistRepo := repositories.NewWishlistRepository()
	addressRepo := repositories.NewAddressRepository()
	orderRepo := repositories.NewOrderRepository()
	couponRepo := repositories.NewCouponRepository()
	adminRepo := repositories.NewAdminRepository()
	reviewRepo := repositories.NewReviewRepository()
	notifRepo := repositories.NewNotificationRepository()
	newsRepo := repositories.NewNewsletterRepository()

	// ==================== SERVICES ====================
	emailService := email.NewService()
	authService := services.NewAuthService(userRepo, emailService)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	wishlistService := services.NewWishlistService(wishlistRepo, productRepo)
	addressService := services.NewAddressService(addressRepo)
	couponService := services.NewCouponService(couponRepo)
	loyaltyService := services.NewLoyaltyService()
	orderService := services.NewOrderService(
		orderRepo, cartRepo, addressRepo, productRepo,
		couponRepo, couponService, loyaltyService, emailService,
	)
	adminService := services.NewAdminService(adminRepo)
	reviewService := services.NewReviewService(reviewRepo, productRepo)
	notifService := services.NewNotificationService(notifRepo)
	newsService := services.NewNewsletterService(newsRepo)

	// ==================== HANDLERS ====================
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	addressHandler := handlers.NewAddressHandler(addressService)
	orderHandler := handlers.NewOrderHandler(orderService)
	adminHandler := handlers.NewAdminHandler(adminService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	notifHandler := handlers.NewNotificationHandler(notifService)
	newsHandler := handlers.NewNewsletterHandler(newsService)

	// ==================== ROUTES ====================
	registerAuthRoutes(api, authHandler)
	registerCategoryRoutes(api, categoryHandler)
	registerProductRoutes(api, productHandler)
	registerCartRoutes(api, cartHandler)
	registerWishlistRoutes(api, wishlistHandler)
	registerAddressRoutes(api, addressHandler)
	registerOrderRoutes(api, orderHandler)
	registerAdminRoutes(api, adminHandler)
	registerReviewRoutes(api, reviewHandler)
	registerNotificationRoutes(api, notifHandler)
	registerNewsletterRoutes(api, newsHandler)
}