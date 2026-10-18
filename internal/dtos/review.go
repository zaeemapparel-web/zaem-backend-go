package dtos

// ==================== REQUESTS ====================

type CreateReviewRequest struct {
	ProductID string   `json:"productId" binding:"required"`
	Rating    int      `json:"rating" binding:"required,min=1,max=5"`
	Title     *string  `json:"title" binding:"omitempty,max=200"`
	Comment   string   `json:"comment" binding:"required,min=5,max=2000"`
	Images    []string `json:"images"`
}

type UpdateReviewRequest struct {
	Rating  *int     `json:"rating" binding:"omitempty,min=1,max=5"`
	Title   *string  `json:"title"`
	Comment *string  `json:"comment" binding:"omitempty,min=5,max=2000"`
	Images  []string `json:"images"`
}

type AdminReviewActionRequest struct {
	IsApproved *bool   `json:"isApproved"`
	AdminReply *string `json:"adminReply"`
}

// ==================== RESPONSES ====================

type ReviewResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userId"`
	UserName      string  `json:"userName"`
	UserAvatar    *string `json:"userAvatar,omitempty"`
	ProductID     string  `json:"productId"`
	Rating        int     `json:"rating"`
	Title         *string `json:"title,omitempty"`
	Comment       string  `json:"comment"`
	Images        []string `json:"images"`
	IsApproved    bool    `json:"isApproved"`
	IsVerified    bool    `json:"isVerified"`
	HelpfulCount  int     `json:"helpfulCount"`
	IsMarkedHelpful bool  `json:"isMarkedHelpful"`
	AdminReply    *string `json:"adminReply,omitempty"`
	AdminRepliedAt *string `json:"adminRepliedAt,omitempty"`
	CreatedAt     string  `json:"createdAt"`
}