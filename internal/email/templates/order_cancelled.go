package templates

import (
	"fmt"
	stdhtml "html"
)

// OrderCancelledEmail — sent when order is cancelled
func OrderCancelledEmail(
	userName, orderNumber, reason, shopURL string,
) (subject, htmlBody string) {
	subject = fmt.Sprintf("Order %s Has Been Cancelled", orderNumber)

	safeName := stdhtml.EscapeString(userName)

	reasonSection := ""
	if reason != "" {
		reasonSection = fmt.Sprintf(`
      <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#FFEBEE;border-left:3px solid #C62828;border-radius:8px;margin:24px 0;">
        <tr>
          <td style="padding:16px 20px;">
            <div style="font-size:11px;letter-spacing:0.2em;color:#C62828;text-transform:uppercase;margin-bottom:6px;">Reason</div>
            <div style="font-size:14px;color:#1D1D1F;line-height:1.6;">%s</div>
          </td>
        </tr>
      </table>
    `, stdhtml.EscapeString(reason))
	}

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      We're sorry to inform you that your order <strong style="color:#1D1D1F;">%s</strong> has been cancelled.
    </p>

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:24px 0 20px;">
      If you were charged, a refund will be processed within 5-7 business days to your original payment method.
    </p>

    %s

    <p style="font-size:13px;color:#86868B;line-height:1.7;margin:24px 0 0;text-align:center;">
      Have questions? Reply to this email or contact us at
      <a href="mailto:zaemlifestyle@gmail.com" style="color:#1D1D1F;text-decoration:underline;">zaemlifestyle@gmail.com</a>
    </p>
  `,
		Heading("Order Cancelled"),
		Paragraph(fmt.Sprintf("Hi %s,", safeName)),
		stdhtml.EscapeString(orderNumber),
		reasonSection,
		Button("Continue Shopping", shopURL),
	)

	preheader := fmt.Sprintf("Order %s has been cancelled", orderNumber)
	htmlBody = BaseLayout(subject, preheader, content)
	return
}