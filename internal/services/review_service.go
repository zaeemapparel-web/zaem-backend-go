package services

import (
	"context"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
)

type ReviewService struct {
	repo        repositories.ReviewRepository
	productRepo repositories.ProductRepository
}

func NewReviewService(r repositories.ReviewRepository, pr repositories.ProductRepository) *ReviewService {
	return &ReviewService{repo: r, productRepo: pr}
}

func (s *ReviewService) GetProductReviews(ctx context.Context, productID string, page, limit int) ([]dtos.ReviewResponse, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 10 }

	reviews, total, err := s.repo.FindByProductID(ctx, productID, page, limit, true)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch reviews")
	}

	result := make([]dtos.ReviewResponse, 0, len(reviews))
	for i := range reviews {
		result = append(result, s.toResponse(&reviews[i]))
	}
	return result, total, nil
}

func (s *ReviewService) Create(ctx context.Context, userID string, req dtos.CreateReviewRequest) (*dtos.ReviewResponse, error) {
	// Check product exists
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil || product == nil {
		return nil, apperrors.ProductNotFound()
	}

	// Check duplicate
	existing, _ := s.repo.FindByUserAndProduct(ctx, userID, req.ProductID)
	if existing != nil {
		return nil, apperrors.Conflict("You have already reviewed this product")
	}

	// Check verified purchase
	isVerified, _ := s.repo.HasPurchased(ctx, userID, req.ProductID)

	review := &models.Review{
		ID:         cuid.New(),
		UserID:     userID,
		ProductID:  req.ProductID,
		Rating:     req.Rating,
		Title:      req.Title,
		Comment:    req.Comment,
		Images:     req.Images,
		IsApproved: false, // Needs admin approval
		IsVerified: isVerified,
	}

	if err := s.repo.Create(ctx, review); err != nil {
		return nil, apperrors.Internal("Failed to create review")
	}

	// Reload to get user details
	review, _ = s.repo.FindByID(ctx, review.ID)
	resp := s.toResponse(review)
	return &resp, nil
}

func (s *ReviewService) Update(ctx context.Context, userID, reviewID string, req dtos.UpdateReviewRequest) (*dtos.ReviewResponse, error) {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil || review == nil {
		return nil, apperrors.NotFound("Review not found")
	}
	if review.UserID != userID {
		return nil, apperrors.Forbidden("You don't own this review")
	}

	updates := map[string]interface{}{}
	if req.Rating != nil { updates["rating"] = *req.Rating }
	if req.Title != nil { updates["title"] = *req.Title }
	if req.Comment != nil { updates["comment"] = *req.Comment }
	if req.Images != nil { updates["images"] = req.Images }
	
	// If edited, need re-approval
	if len(updates) > 0 {
		updates["is_approved"] = false
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateFields(ctx, reviewID, updates); err != nil {
			return nil, apperrors.Internal("Failed to update review")
		}
	}

	review, _ = s.repo.FindByID(ctx, reviewID)
	resp := s.toResponse(review)
	return &resp, nil
}

func (s *ReviewService) Delete(ctx context.Context, userID string, isAdmin bool, reviewID string) error {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil || review == nil {
		return apperrors.NotFound("Review not found")
	}
	if !isAdmin && review.UserID != userID {
		return apperrors.Forbidden("You don't own this review")
	}

	if err := s.repo.Delete(ctx, reviewID); err != nil {
		return apperrors.Internal("Failed to delete review")
	}
	return nil
}

func (s *ReviewService) MarkHelpful(ctx context.Context, userID, reviewID string) error {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil || review == nil {
		return apperrors.NotFound("Review not found")
	}
	if review.UserID == userID {
		return apperrors.BadRequest("You cannot mark your own review as helpful")
	}

	existing, _ := s.repo.FindHelpful(ctx, reviewID, userID)
	
	if existing != nil {
		// Toggle off
		_ = s.repo.DeleteHelpful(ctx, existing.ID)
		_ = s.repo.IncrementHelpfulCount(ctx, reviewID, -1)
		return nil
	}

	// Add new
	h := &models.ReviewHelpful{
		ID:       cuid.New(),
		ReviewID: reviewID,
		UserID:   userID,
	}
	if err := s.repo.CreateHelpful(ctx, h); err != nil {
		return apperrors.Internal("Failed to mark helpful")
	}
	_ = s.repo.IncrementHelpfulCount(ctx, reviewID, 1)
	return nil
}

// ==================== ADMIN ====================

func (s *ReviewService) GetAllAdmin(ctx context.Context, page, limit int, isApproved *bool) ([]dtos.ReviewResponse, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }

	reviews, total, err := s.repo.FindAll(ctx, page, limit, isApproved)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch reviews")
	}

	result := make([]dtos.ReviewResponse, 0, len(reviews))
	for i := range reviews {
		result = append(result, s.toResponse(&reviews[i]))
	}
	return result, total, nil
}

func (s *ReviewService) AdminAction(ctx context.Context, reviewID string, req dtos.AdminReviewActionRequest) error {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil || review == nil {
		return apperrors.NotFound("Review not found")
	}

	updates := map[string]interface{}{}
	if req.IsApproved != nil {
		updates["is_approved"] = *req.IsApproved
	}
	if req.AdminReply != nil {
		updates["admin_reply"] = *req.AdminReply
		updates["admin_replied_at"] = time.Now().UTC()
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateFields(ctx, reviewID, updates); err != nil {
			return apperrors.Internal("Failed to update review")
		}
	}
	return nil
}

// ==================== HELPERS ====================

func (s *ReviewService) toResponse(r *models.Review) dtos.ReviewResponse {
	resp := dtos.ReviewResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		ProductID:    r.ProductID,
		Rating:       r.Rating,
		Title:        r.Title,
		Comment:      r.Comment,
		Images:       r.Images,
		IsApproved:   r.IsApproved,
		IsVerified:   r.IsVerified,
		HelpfulCount: r.HelpfulCount,
		AdminReply:   r.AdminReply,
		CreatedAt:    r.CreatedAt.Format(time.RFC3339),
	}

	if r.User.ID != "" {
		resp.UserName = r.User.Name
		resp.UserAvatar = r.User.Avatar
	}

	if r.AdminRepliedAt != nil {
		t := r.AdminRepliedAt.Format(time.RFC3339)
		resp.AdminRepliedAt = &t
	}

	return resp
}
// ==================== GET MY REVIEWS ====================

func (s *ReviewService) GetMyReviews(ctx context.Context, userID string) ([]dtos.ReviewResponse, error) {
	// Repository mein FindByUserID ka function chahiye hoga.
	// Filhal hum khali list return kar rahe hain taake build pass ho jaye.
	// (Baad mein isay database se connect karenge)
	return []dtos.ReviewResponse{}, nil
}

// ==================== GET REVIEW STATS ====================

func (s *ReviewService) GetReviewStats(ctx context.Context) (map[string]interface{}, error) {
	// Admin panel ke liye review stats
	return map[string]interface{}{
		"totalReviews":    0,
		"pendingReviews":  0,
		"approvedReviews": 0,
		"averageRating":   0,
	}, nil
}