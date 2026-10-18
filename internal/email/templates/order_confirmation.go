package templates

import (
	"fmt"
	stdhtml "html"
	"strings"
)

// OrderItem — item data for template
type OrderItem struct {
	Name     string
	Quantity int
	Price    float64
	Size     string
	Color    string
	Image    string
}

// OrderConfirmationEmail — sent when order is placed
func OrderConfirmationEmail(
	userName, orderNumber, orderDate string,
	total float64,
	items []OrderItem,
	orderURL string,
) (subject, htmlBody string) {
	subject = fmt.Sprintf("Order Confirmed - %s", orderNumber)

	safeName := stdhtml.EscapeString(userName)

	// Build items HTML
	var itemsHTML strings.Builder
	for _, item := range items {
		variant := ""
		if item.Size != "" {
			variant += fmt.Sprintf("Size: %s ", stdhtml.EscapeString(item.Size))
		}
		if item.Color != "" {
			variant += fmt.Sprintf("| Color: %s", stdhtml.EscapeString(item.Color))
		}

		itemsHTML.WriteString(fmt.Sprintf(`
      <tr>
        <td style="padding:16px 0;border-bottom:1px solid #F0F0F0;">
          <div style="font-size:14px;color:#1D1D1F;font-weight:500;">%s</div>
          <div style="font-size:12px;color:#86868B;margin-top:4px;">%s</div>
          <div style="font-size:12px;color:#86868B;margin-top:4px;">Qty: %d</div>
        </td>
        <td style="padding:16px 0;border-bottom:1px solid #F0F0F0;text-align:right;vertical-align:top;">
          <div style="font-size:14px;color:#1D1D1F;font-weight:500;">Rs. %.0f</div>
        </td>
      </tr>
    `, stdhtml.EscapeString(item.Name), variant, item.Quantity, item.Price*float64(item.Quantity)))
	}

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      Thank you for your order! We're getting it ready and will notify you once it ships.
    </p>

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#FAFAFA;border-radius:12px;margin:24px 0;">
      <tr>
        <td style="padding:20px;">
          <table width="100%%" cellspacing="0" cellpadding="0" border="0">
            <tr>
              <td style="font-size:11px;letter-spacing:0.2em;color:#86868B;text-transform:uppercase;padding-bottom:4px;">Order Number</td>
              <td style="font-size:11px;letter-spacing:0.2em;color:#86868B;text-transform:uppercase;padding-bottom:4px;text-align:right;">Order Date</td>
            </tr>
            <tr>
              <td style="font-size:15px;color:#1D1D1F;font-weight:600;">%s</td>
              <td style="font-size:15px;color:#1D1D1F;text-align:right;">%s</td>
            </tr>
          </table>
        </td>
      </tr>
    </table>

    <div style="font-size:11px;letter-spacing:0.2em;color:#86868B;text-transform:uppercase;margin:32px 0 8px;">Order Summary</div>

    <table width="100%%" cellspacing="0" cellpadding="0" border="0">
      %s
      <tr>
        <td style="padding:20px 0 0;font-size:14px;color:#1D1D1F;font-weight:600;">Total</td>
        <td style="padding:20px 0 0;text-align:right;font-size:20px;color:#1D1D1F;font-family:Georgia,serif;">Rs. %.0f</td>
      </tr>
    </table>

    %s

    <p style="font-size:13px;color:#86868B;line-height:1.7;margin:24px 0 0;">
      You'll receive an email when your order ships with tracking information.
    </p>
  `,
		Heading("Order Confirmed"),
		Paragraph(fmt.Sprintf("Hi %s,", safeName)),
		stdhtml.EscapeString(orderNumber),
		orderDate,
		itemsHTML.String(),
		total,
		Button("View Order", orderURL),
	)

	preheader := fmt.Sprintf("Order %s confirmed — we're preparing it now", orderNumber)
	htmlBody = BaseLayout(subject, preheader, content)
	return
}