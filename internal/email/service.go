package email

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/config"
	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/email/templates"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/pkg/cuid"
)

type Service struct {
	sender *ResendSender
}

var defaultService *Service

func NewService() *Service {
	s := &Service{sender: NewResendSender()}
	defaultService = s
	return s
}

// ==================== AUTH EMAILS ====================

func (s *Service) SendWelcome(user *models.User) {
	loginURL := fmt.Sprintf("%s/account", config.C.FrontendURL)
	shopURL := fmt.Sprintf("%s/shop", config.C.FrontendURL)

	subject, html := templates.WelcomeEmail(user.Name, loginURL, shopURL)
	s.sendAsync(&SendRequest{
		To: user.Email, Subject: subject, HTML: html, Type: EmailWelcome, UserID: &user.ID,
	})
}

func (s *Service) SendPasswordReset(user *models.User, token string, expiryMinutes int) {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.C.FrontendURL, token)

	subject, html := templates.PasswordResetEmail(user.Name, resetURL, expiryMinutes)
	s.sendAsync(&SendRequest{
		To: user.Email, Subject: subject, HTML: html, Type: EmailPasswordReset, UserID: &user.ID,
	})
}

func (s *Service) SendPasswordResetSuccess(user *models.User) {
	loginURL := fmt.Sprintf("%s/account/login", config.C.FrontendURL)

	subject, html := templates.PasswordResetSuccessEmail(user.Name, loginURL)
	s.sendAsync(&SendRequest{
		To: user.Email, Subject: subject, HTML: html, Type: EmailResetSuccess, UserID: &user.ID,
	})
}

// ==================== ORDER EMAILS ====================

func (s *Service) SendOrderConfirmation(order *models.Order) {
	if order.User.Email == "" {
		return
	}

	// Build items
	items := make([]templates.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		size := ""
		if item.Size != nil {
			size = *item.Size
		}
		color := ""
		if item.Color != nil {
			color = *item.Color
		}
		image := ""
		if item.Image != nil {
			image = *item.Image
		}

		items = append(items, templates.OrderItem{
			Name: item.Name, Quantity: item.Quantity, Price: item.Price,
			Size: size, Color: color, Image: image,
		})
	}

	orderURL := fmt.Sprintf("%s/account/orders/%s", config.C.FrontendURL, order.ID)
	orderDate := order.CreatedAt.Format("Jan 02, 2006")

	subject, html := templates.OrderConfirmationEmail(
		order.User.Name, order.OrderNumber, orderDate, order.Total, items, orderURL,
	)

	orderID := order.ID
	s.sendAsync(&SendRequest{
		To: order.User.Email, Subject: subject, HTML: html,
		Type: EmailOrderConfirm, UserID: &order.UserID, OrderID: &orderID,
	})
}

func (s *Service) SendOrderShipped(order *models.Order) {
	if order.User.Email == "" {
		return
	}

	trackingNumber := ""
	if order.TrackingNumber != nil {
		trackingNumber = *order.TrackingNumber
	}
	trackingURL := ""
	if order.TrackingURL != nil {
		trackingURL = *order.TrackingURL
	}

	orderURL := fmt.Sprintf("%s/account/orders/%s", config.C.FrontendURL, order.ID)

	subject, html := templates.OrderShippedEmail(
		order.User.Name, order.OrderNumber, trackingNumber, trackingURL, orderURL,
	)

	orderID := order.ID
	s.sendAsync(&SendRequest{
		To: order.User.Email, Subject: subject, HTML: html,
		Type: EmailOrderShipped, UserID: &order.UserID, OrderID: &orderID,
	})
}

func (s *Service) SendOrderDelivered(order *models.Order) {
	if order.User.Email == "" {
		return
	}

	reviewURL := fmt.Sprintf("%s/account/orders/%s", config.C.FrontendURL, order.ID)
	shopURL := fmt.Sprintf("%s/shop", config.C.FrontendURL)

	subject, html := templates.OrderDeliveredEmail(
		order.User.Name, order.OrderNumber, reviewURL, shopURL,
	)

	orderID := order.ID
	s.sendAsync(&SendRequest{
		To: order.User.Email, Subject: subject, HTML: html,
		Type: EmailOrderDelivered, UserID: &order.UserID, OrderID: &orderID,
	})
}

func (s *Service) SendOrderCancelled(order *models.Order, reason string) {
	if order.User.Email == "" {
		return
	}

	shopURL := fmt.Sprintf("%s/shop", config.C.FrontendURL)

	subject, html := templates.OrderCancelledEmail(
		order.User.Name, order.OrderNumber, reason, shopURL,
	)

	orderID := order.ID
	s.sendAsync(&SendRequest{
		To: order.User.Email, Subject: subject, HTML: html,
		Type: EmailOrderCancelled, UserID: &order.UserID, OrderID: &orderID,
	})
}

// ==================== INTERNAL ====================

func (s *Service) sendAsync(req *SendRequest) {
	go func() {
		emailLog := &models.EmailLog{
			ID: cuid.New(), To: req.To, Subject: req.Subject,
			Template: string(req.Type), OrderID: req.OrderID, UserID: req.UserID,
			Status: "pending",
		}
		database.DB.Create(emailLog)

		resendID, err := s.sender.Send(req.To, req.Subject, req.HTML)

		if err != nil {
			log.Error().Err(err).Str("to", req.To).Str("template", string(req.Type)).Msg("❌ Failed to send email")
			errStr := err.Error()
			database.DB.Model(emailLog).Updates(map[string]interface{}{
				"status": "failed", "error": errStr,
			})
			return
		}

		log.Info().Str("to", req.To).Str("template", string(req.Type)).Str("resendId", resendID).Msg("✅ Email sent")

		now := time.Now().UTC()
		database.DB.Model(emailLog).Updates(map[string]interface{}{
			"status": "sent", "resend_id": resendID, "sent_at": now,
		})
	}()
}

func GetDefault() *Service {
	return defaultService
}