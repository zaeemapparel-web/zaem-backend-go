package services

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/email"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
	"github.com/zaem/backend-go/pkg/hash"
	"github.com/zaem/backend-go/pkg/jwt"
	"github.com/zaem/backend-go/pkg/random"
)

const (
	ResetTokenExpiryMinutes = 30
	EmailVerifyExpiryHours  = 24
)

// AuthService handles all auth business logic
type AuthService struct {
	userRepo repositories.UserRepository
	emailSvc *email.Service
}

// NewAuthService creates new instance
func NewAuthService(userRepo repositories.UserRepository, emailSvc *email.Service) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		emailSvc: emailSvc,
	}
}

// ==================== REGISTER ====================

func (s *AuthService) Register(ctx context.Context, req dtos.RegisterRequest) (*dtos.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check email exists
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check email existence")
		return nil, apperrors.Internal("Failed to register user")
	}
	if exists {
		return nil, apperrors.EmailAlreadyExists()
	}

	// Check phone exists (if provided)
	var phonePtr *string
	if req.Phone != "" {
		exists, err = s.userRepo.ExistsByPhone(ctx, req.Phone)
		if err != nil {
			return nil, apperrors.Internal("Failed to register user")
		}
		if exists {
			return nil, apperrors.PhoneAlreadyExists()
		}
		phone := strings.TrimSpace(req.Phone)
		phonePtr = &phone
	}

	// Hash password
	hashedPassword, err := hash.Password(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, apperrors.Internal("Failed to register user")
	}

	// Generate unique referral code
	referralCode := strings.ToUpper(random.Token(4))

	// Generate email verify token
	verifyToken := random.Token(32)

	// Create user
	user := &models.User{
		ID:               cuid.New(),
		Name:             strings.TrimSpace(req.Name),
		Email:            email,
		Phone:            phonePtr,
		Password:         hashedPassword,
		Role:             models.RoleUser,
		IsActive:         true,
		IsEmailVerified:  false,
		EmailVerifyToken: &verifyToken,
		Tier:             models.TierBronze,
		ReferralCode:     &referralCode,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, apperrors.Internal("Failed to register user")
	}

	// Send welcome email (async, non-blocking)
	go s.emailSvc.SendWelcome(user)

	log.Info().
		Str("userId", user.ID).
		Str("email", user.Email).
		Msg("✅ User registered")

	// Generate JWT
	token, err := jwt.Generate(user.ID, user.Email, string(user.Role))
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate JWT")
		return nil, apperrors.Internal("Failed to register user")
	}

	return &dtos.AuthResponse{
		User:  toUserResponse(user),
		Token: token,
	}, nil
}

// ==================== LOGIN ====================

func (s *AuthService) Login(ctx context.Context, req dtos.LoginRequest, ip string) (*dtos.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user")
		return nil, apperrors.Internal("Failed to login")
	}
	if user == nil {
		return nil, apperrors.InvalidCredentials()
	}

	// Check password
	if !hash.Compare(user.Password, req.Password) {
		return nil, apperrors.InvalidCredentials()
	}

	// Check active
	if !user.IsActive {
		return nil, apperrors.AccountDisabled()
	}

	// Update last login
	now := time.Now().UTC()
	_ = s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	})

	// Generate JWT
	token, err := jwt.Generate(user.ID, user.Email, string(user.Role))
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate JWT")
		return nil, apperrors.Internal("Failed to login")
	}

	log.Info().
		Str("userId", user.ID).
		Str("email", user.Email).
		Msg("✅ User logged in")

	return &dtos.AuthResponse{
		User:  toUserResponse(user),
		Token: token,
	}, nil
}

// ==================== GET ME ====================

func (s *AuthService) GetMe(ctx context.Context, userID string) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch user")
	}
	if user == nil {
		return nil, apperrors.UserNotFound()
	}

	resp := toUserResponse(user)
	return &resp, nil
}

// ==================== UPDATE PROFILE ====================

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req dtos.UpdateProfileRequest) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, apperrors.UserNotFound()
	}

	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}

	if req.Phone != "" {
		// Check if phone belongs to another user
		existing, _ := s.userRepo.FindByPhone(ctx, req.Phone)
		if existing != nil && existing.ID != userID {
			return nil, apperrors.PhoneAlreadyExists()
		}
		updates["phone"] = strings.TrimSpace(req.Phone)
	}

	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		resp := toUserResponse(user)
		return &resp, nil
	}

	if err := s.userRepo.UpdateFields(ctx, userID, updates); err != nil {
		return nil, apperrors.Internal("Failed to update profile")
	}

	// Reload
	user, _ = s.userRepo.FindByID(ctx, userID)
	resp := toUserResponse(user)
	return &resp, nil
}

// ==================== CHANGE PASSWORD ====================

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req dtos.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return apperrors.UserNotFound()
	}

	if !hash.Compare(user.Password, req.CurrentPassword) {
		return apperrors.BadRequest("Current password is incorrect")
	}

	newHash, err := hash.Password(req.NewPassword)
	if err != nil {
		return apperrors.Internal("Failed to change password")
	}

	if err := s.userRepo.UpdateFields(ctx, userID, map[string]interface{}{
		"password": newHash,
	}); err != nil {
		return apperrors.Internal("Failed to change password")
	}

	return nil
}

// ==================== FORGOT PASSWORD ====================

func (s *AuthService) ForgotPassword(ctx context.Context, req dtos.ForgotPasswordRequest) error {
	emailAddr := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.FindByEmail(ctx, emailAddr)
	if err != nil {
		return apperrors.Internal("Failed to process request")
	}

	// SECURITY: Always return success even if user doesn't exist
	// (prevents email enumeration attacks)
	if user == nil {
		log.Info().Str("email", emailAddr).Msg("⚠️  Forgot password for non-existent email")
		return nil
	}

	// Generate reset token
	token := random.Token(32)
	expiry := time.Now().UTC().Add(ResetTokenExpiryMinutes * time.Minute)

	if err := s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"reset_token":        token,
		"reset_token_expiry": expiry,
	}); err != nil {
		return apperrors.Internal("Failed to process request")
	}

	// Send email (async)
	go s.emailSvc.SendPasswordReset(user, token, ResetTokenExpiryMinutes)

	log.Info().
		Str("userId", user.ID).
		Str("email", user.Email).
		Msg("📧 Password reset email sent")

	return nil
}

// ==================== RESET PASSWORD ====================

func (s *AuthService) ResetPassword(ctx context.Context, req dtos.ResetPasswordRequest) error {
	user, err := s.userRepo.FindByResetToken(ctx, req.Token)
	if err != nil {
		return apperrors.Internal("Failed to reset password")
	}
	if user == nil {
		return apperrors.InvalidResetToken()
	}

	// Check expiry
	if user.ResetTokenExpiry == nil || user.ResetTokenExpiry.Before(time.Now().UTC()) {
		return apperrors.ResetTokenExpired()
	}

	// Hash new password
	newHash, err := hash.Password(req.Password)
	if err != nil {
		return apperrors.Internal("Failed to reset password")
	}

	// Update password + clear token (single-use)
	if err := s.userRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"password":           newHash,
		"reset_token":        nil,
		"reset_token_expiry": nil,
	}); err != nil {
		return apperrors.Internal("Failed to reset password")
	}

	// Send confirmation email
	go s.emailSvc.SendPasswordResetSuccess(user)

	log.Info().
		Str("userId", user.ID).
		Msg("✅ Password reset successful")

	return nil
}

// ==================== HELPER ====================

func toUserResponse(u *models.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Phone:           u.Phone,
		Role:            string(u.Role),
		IsEmailVerified: u.IsEmailVerified,
		LoyaltyPoints:   u.LoyaltyPoints,
		Tier:            string(u.Tier),
		ReferralCode:    u.ReferralCode,
		Avatar:          u.Avatar,
	}
}