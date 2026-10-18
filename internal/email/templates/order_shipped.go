package templates

import (
	"fmt"
	stdhtml "html"
)

// OrderShippedEmail — sent when order is shipped
func OrderShippedEmail(
	userName, orderNumber, trackingNumber, trackingURL, orderURL string,
) (subject, htmlBody string) {
	subject = fmt.Sprintf("Your Order %s Has Shipped", orderNumber)

	safeName := stdhtml.EscapeString(userName)

	trackingSection := ""
	if trackingNumber != "" {
		trackingButton := Button("Track Order", trackingURL)
		trackingSection = fmt.Sprintf(`
      <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#E8F5E9;border-radius:12px;margin:24px 0;">
        <tr>
          <td style="padding:20px;">
            <div style="font-size:11px;letter-spacing:0.2em;color:#2E7D32;text-transform:uppercase;margin-bottom:8px;">Tracking Number</div>
            <div style="font-size:16px;color:#1D1D1F;font-weight:600;font-family:monospace;">%s</div>
          </td>
        </tr>
      </table>
      %s
    `, stdhtml.EscapeString(trackingNumber), trackingButton)
	} else {
		trackingSection = Button("View Order", orderURL)
	}

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      Great news! Your order <strong style="color:#1D1D1F;">%s</strong> has been shipped and is on its way to you.
    </p>

    %s

    <p style="font-size:13px;color:#86868B;line-height:1.7;margin:24px 0 0;">
      You can track your order using the link above. Delivery usually takes 2-4 business days.
    </p>
  `,
		Heading("Your Order Has Shipped 🚚"),
		Paragraph(fmt.Sprintf("Hi %s,", safeName)),
		stdhtml.EscapeString(orderNumber),
		trackingSection,
	)

	preheader := fmt.Sprintf("Order %s is on its way to you", orderNumber)
	htmlBody = BaseLayout(subject, preheader, content)
	return
}