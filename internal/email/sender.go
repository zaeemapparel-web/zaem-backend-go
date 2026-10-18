package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/config"
)

// Resend API endpoint
const resendAPIURL = "https://api.resend.com/emails"

// resendRequest — payload for Resend API
type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	ReplyTo string   `json:"reply_to,omitempty"`
}

// resendResponse — response from Resend API
type resendResponse struct {
	ID string `json:"id"`
}

// ResendSender — handles sending emails via Resend
type ResendSender struct {
	apiKey    string
	from      string
	replyTo   string
	client    *http.Client
}

// NewResendSender creates a new sender
func NewResendSender() *ResendSender {
	return &ResendSender{
		apiKey:  config.C.ResendAPIKey,
		from:    config.C.EmailFrom,
		replyTo: config.C.EmailReplyTo,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Send sends an email via Resend API
// Returns the Resend message ID (or empty if failed) and error
func (s *ResendSender) Send(to, subject, html string) (string, error) {
	if s.apiKey == "" {
		log.Warn().Msg("⚠️  Resend API key not set — email not sent")
		return "", fmt.Errorf("resend api key not configured")
	}

	if s.from == "" {
		return "", fmt.Errorf("email from address not configured")
	}

	payload := resendRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
		ReplyTo: s.replyTo,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", resendAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("resend API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result resendResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.ID, nil
}