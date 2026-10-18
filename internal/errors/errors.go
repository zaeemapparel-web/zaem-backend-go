package errors

import "net/http"

// AppError is our custom error type
type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// ==================== CONSTRUCTORS ====================

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func Wrap(err error, status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Err: err}
}

// ==================== GENERIC ====================

func BadRequest(msg string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", msg)
}

func Unauthorized(msg string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

func Forbidden(msg string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", msg)
}

func NotFound(msg string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", msg)
}

func Conflict(msg string) *AppError {
	return New(http.StatusConflict, "CONFLICT", msg)
}

func Validation(msg string) *AppError {
	return New(http.StatusUnprocessableEntity, "VALIDATION_ERROR", msg)
}

func TooManyRequests(msg string) *AppError {
	return New(http.StatusTooManyRequests, "RATE_LIMIT", msg)
}

func Internal(msg string) *AppError {
	if msg == "" {
		msg = "Internal server error"
	}
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

// ==================== AUTH ====================

func InvalidCredentials() *AppError {
	return New(http.StatusUnauthorized, "AUTH_001", "Invalid email or password")
}

func EmailAlreadyExists() *AppError {
	return New(http.StatusConflict, "AUTH_002", "Email already registered")
}

func PhoneAlreadyExists() *AppError {
	return New(http.StatusConflict, "AUTH_003", "Phone number already registered")
}

func AccountDisabled() *AppError {
	return New(http.StatusForbidden, "AUTH_004", "Your account has been disabled")
}

func InvalidToken() *AppError {
	return New(http.StatusUnauthorized, "AUTH_005", "Invalid or expired token")
}

func ResetTokenExpired() *AppError {
	return New(http.StatusBadRequest, "AUTH_006", "Password reset link has expired")
}

func InvalidResetToken() *AppError {
	return New(http.StatusBadRequest, "AUTH_007", "Invalid password reset token")
}

func UserNotFound() *AppError {
	return New(http.StatusNotFound, "AUTH_008", "User not found")
}

func WeakPassword() *AppError {
	return New(http.StatusBadRequest, "AUTH_009", "Password is too weak")
}

// ==================== RESOURCE ====================

func CategoryNotFound() *AppError {
	return New(http.StatusNotFound, "CAT_001", "Category not found")
}

func CategorySlugExists() *AppError {
	return New(http.StatusConflict, "CAT_002", "Category with this slug already exists")
}

func CategoryHasProducts() *AppError {
	return New(http.StatusBadRequest, "CAT_003", "Cannot delete category with products")
}

func ProductNotFound() *AppError {
	return New(http.StatusNotFound, "PROD_001", "Product not found")
}

func ProductSlugExists() *AppError {
	return New(http.StatusConflict, "PROD_002", "Product with this slug already exists")
}

// ==================== CART ====================

func CartItemNotFound() *AppError {
	return New(http.StatusNotFound, "CART_001", "Cart item not found")
}

func ProductOutOfStock() *AppError {
	return New(http.StatusBadRequest, "CART_002", "Product is out of stock")
}

func InsufficientStock(available int) *AppError {
	return New(http.StatusBadRequest, "CART_003", "Insufficient stock available")
}

// ==================== ORDER ====================

func OrderNotFound() *AppError {
	return New(http.StatusNotFound, "ORD_001", "Order not found")
}

func EmptyCart() *AppError {
	return New(http.StatusBadRequest, "ORD_002", "Your cart is empty")
}

func OrderNotCancellable() *AppError {
	return New(http.StatusBadRequest, "ORD_003", "Order cannot be cancelled at this stage")
}