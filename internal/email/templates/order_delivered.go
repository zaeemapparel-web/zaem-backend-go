package templates

import (
	"fmt"
	stdhtml "html"
)

// OrderDeliveredEmail — sent when order is delivered
func OrderDeliveredEmail(
	userName, orderNumber, reviewURL, shopURL string,
) (subject, htmlBody string) {
	subject = fmt.Sprintf("Your Order %s Has Been Delivered", orderNumber)

	safeName := stdhtml.EscapeString(userName)

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      Your order <strong style="color:#1D1D1F;">%s</strong> has been delivered successfully.
      We hope you love it!
    </p>

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#E8F5E9;border-left:3px solid #2E7D32;border-radius:8px;margin:24px 0;">
      <tr>
        <td style="padding:16px 20px;">
          <div style="font-size:13px;color:#2E7D32;line-height:1.6;">
            <strong>✓ Delivered Successfully</strong>
          </div>
        </td>
      </tr>
    </table>

    <div style="font-size:11px;letter-spacing:0.2em;color:#86868B;text-transform:uppercase;margin:32px 0 12px;text-align:center;">How was your experience?</div>

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;text-align:center;">
      Share your thoughts and help others make better choices.
      Your review means a lot to us!
    </p>

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:0 auto;">
      <tr>
        <td style="background-color:#1D1D1F;border-radius:12px;">
          <a href="%s" style="display:inline-block;padding:16px 40px;color:#FFFFFF;text-decoration:none;font-size:12px;letter-spacing:0.2em;text-transform:uppercase;font-weight:500;">Write a Review</a>
        </td>
      </tr>
    </table>

    <div style="margin-top:24px;text-align:center;">
      <a href="%s" style="color:#6E6E73;font-size:12px;text-decoration:underline;">Continue Shopping</a>
    </div>
  `,
		Heading("Order Delivered ✓"),
		Paragraph(fmt.Sprintf("Hi %s,", safeName)),
		stdhtml.EscapeString(orderNumber),
		reviewURL,
		shopURL,
	)

	preheader := fmt.Sprintf("Order %s delivered — how did we do?", orderNumber)
	htmlBody = BaseLayout(subject, preheader, content)
	return
}