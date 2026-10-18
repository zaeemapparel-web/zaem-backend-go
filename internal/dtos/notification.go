package dtos

type NotificationResponse struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Title    string  `json:"title"`
	Message  string  `json:"message"`
	Link     *string `json:"link,omitempty"`
	ImageURL *string `json:"imageUrl,omitempty"`
	IsRead   bool    `json:"isRead"`
	CreatedAt string `json:"createdAt"`
}