package models

// Role
type Role string

const (
	RoleUser    Role = "USER"
	RoleAdmin   Role = "ADMIN"
	RoleManager Role = "MANAGER"
	RoleStaff   Role = "STAFF"
)

// OrderStatus
type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderConfirmed  OrderStatus = "CONFIRMED"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderShipped    OrderStatus = "SHIPPED"
	OrderDelivered  OrderStatus = "DELIVERED"
	OrderCancelled  OrderStatus = "CANCELLED"
	OrderReturned   OrderStatus = "RETURNED"
)

// PaymentMethod
type PaymentMethod string

const (
	PaymentCOD          PaymentMethod = "COD"
	PaymentJazzCash     PaymentMethod = "JAZZCASH"
	PaymentEasypaisa    PaymentMethod = "EASYPAISA"
	PaymentCard         PaymentMethod = "CARD"
	PaymentBankTransfer PaymentMethod = "BANK_TRANSFER"
)

// PaymentStatus
type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "PENDING"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentFailed   PaymentStatus = "FAILED"
	PaymentRefunded PaymentStatus = "REFUNDED"
)

// CouponType
type CouponType string

const (
	CouponPercentage   CouponType = "PERCENTAGE"
	CouponFixed        CouponType = "FIXED"
	CouponFreeShipping CouponType = "FREE_SHIPPING"
	CouponBuyXGetY     CouponType = "BUY_X_GET_Y"
)

// LoyaltyTier
type LoyaltyTier string

const (
	TierBronze   LoyaltyTier = "BRONZE"
	TierSilver   LoyaltyTier = "SILVER"
	TierGold     LoyaltyTier = "GOLD"
	TierPlatinum LoyaltyTier = "PLATINUM"
)

// LoyaltyAction
type LoyaltyAction string

const (
	LoyaltyOrderPlaced   LoyaltyAction = "ORDER_PLACED"
	LoyaltyReviewWritten LoyaltyAction = "REVIEW_WRITTEN"
	LoyaltyPhotoReview   LoyaltyAction = "PHOTO_REVIEW"
	LoyaltyReferral      LoyaltyAction = "REFERRAL"
	LoyaltySignup        LoyaltyAction = "SIGNUP"
	LoyaltyBirthday      LoyaltyAction = "BIRTHDAY"
	LoyaltyRedeemed      LoyaltyAction = "REDEEMED"
	LoyaltyAdminAdjust   LoyaltyAction = "ADMIN_ADJUSTMENT"
)

// NotificationType
type NotificationType string

const (
	NotifOrderUpdate NotificationType = "ORDER_UPDATE"
	NotifPriceDrop   NotificationType = "PRICE_DROP"
	NotifBackInStock NotificationType = "BACK_IN_STOCK"
	NotifPromotion   NotificationType = "PROMOTION"
	NotifLoyalty     NotificationType = "LOYALTY"
	NotifReviewReply NotificationType = "REVIEW_REPLY"
)