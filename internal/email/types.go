package email

import "github.com/zaem/backend-go/internal/models"

// EmailType — different types of emails
type EmailType string

const (
	EmailWelcome        EmailType = "welcome"
	EmailPasswordReset  EmailType = "password_reset"
	EmailResetSuccess   EmailType = "reset_success"
	EmailVerifyEmail    EmailType = "verify_email"
	EmailOrderConfirm   EmailType = "order_confirmation"
	EmailOrderShipped   EmailType = "order_shipped"
	EmailOrderDelivered EmailType = "order_delivered"
	EmailOrderCancelled EmailType = "order_cancelled"
	EmailOrderReturned  EmailType = "order_returned"
)

type SendRequest struct {
	To      string
	Subject string
	HTML    string
	Type    EmailType
	UserID  *string
	OrderID *string
}

func userFullName(u *models.User) string {
	if u == nil {
		return "Customer"
	}
	return u.Name
}