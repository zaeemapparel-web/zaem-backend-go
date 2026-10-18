package templates

import (
	"fmt"
	stdhtml "html"
)

// WelcomeEmail — sent when user registers
func WelcomeEmail(userName, loginURL, shopURL string) (subject, htmlBody string) {
	subject = "Welcome to ZAEM — Style. Redefined."

	safeName := stdhtml.EscapeString(userName)

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      We're thrilled to welcome you to <strong style="color:#1D1D1F;">ZAEM</strong> — where style meets substance.
    </p>

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      Your account is now active. Explore our curated collection of premium clothing, signature fragrances, and artisan-crafted bags.
    </p>

    %s

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#FAFAFA;border-radius:12px;margin:24px 0;">
      <tr>
        <td style="padding:24px;">
          <div style="font-size:11px;letter-spacing:0.2em;color:#86868B;text-transform:uppercase;margin-bottom:12px;">
            Your Member Benefits
          </div>
          <div style="font-size:14px;color:#1D1D1F;line-height:2;">
            ✦ &nbsp;Early access to new collections<br>
            ✦ &nbsp;Exclusive member-only offers<br>
            ✦ &nbsp;Loyalty points on every order<br>
            ✦ &nbsp;Free shipping above Rs. 5,000
          </div>
        </td>
      </tr>
    </table>

    %s

    %s

    <p style="font-size:14px;color:#86868B;line-height:1.7;margin:24px 0 0;text-align:center;">
      Follow us on <a href="https://instagram.com/zaemstore" style="color:#1D1D1F;text-decoration:underline;">Instagram</a> for the latest drops.
    </p>
  `,
		Heading(fmt.Sprintf("Welcome, %s", safeName)),
		Paragraph("Thank you for joining the ZAEM family."),
		Button("Start Shopping", shopURL),
		Divider(),
		Paragraph(fmt.Sprintf("Need help getting started? <a href=\"%s\" style=\"color:#1D1D1F;text-decoration:underline;\">Visit your account</a> to complete your profile.", loginURL)),
	)

	preheader := "Welcome to ZAEM — your journey starts here."
	htmlBody = BaseLayout(subject, preheader, content)
	return
}