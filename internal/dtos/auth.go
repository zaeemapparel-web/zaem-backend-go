package dtos

// ==================== REQUESTS ====================

// RegisterRequest — matches frontend's register form
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=15"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// LoginRequest — matches frontend's login form
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest — matches frontend's profile update
type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Phone string `json:"phone" binding:"omitempty,min=10,max=15"`
	Avatar string `json:"avatar" binding:"omitempty,url"`
}

// ChangePasswordRequest
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=100"`
}

// ForgotPasswordRequest — matches frontend
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest — matches frontend
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// VerifyEmailRequest — for email verification
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshTokenRequest — for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ==================== RESPONSES ====================

// UserResponse — EXACTLY what frontend expects
// Frontend stores this in localStorage after login/register
type UserResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Phone           *string `json:"phone,omitempty"`
	Role            string  `json:"role"`
	IsEmailVerified bool    `json:"isEmailVerified"`
	LoyaltyPoints   int     `json:"loyaltyPoints"`
	Tier            string  `json:"tier"`
	ReferralCode    *string `json:"referralCode,omitempty"`
	Avatar          *string `json:"avatar,omitempty"`
}

// AuthResponse — what frontend expects after register/login
// { success: true, data: { user, token } }
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// FullAuthResponse — with refresh token (for internal use)
type FullAuthResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken,omitempty"`
}