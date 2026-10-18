package dtos

type NewsletterSubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type NewsletterUnsubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}