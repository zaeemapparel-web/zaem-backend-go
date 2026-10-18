package templates

import (
	"fmt"
	stdhtml "html"
)

// PasswordResetEmail — sent when user requests password reset
func PasswordResetEmail(userName, resetURL string, expiryMinutes int) (subject, htmlBody string) {
	subject = "Reset Your ZAEM Password"

	safeName := stdhtml.EscapeString(userName)

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      We received a request to reset the password for your ZAEM account.
      Click the button below to create a new password.
    </p>

    %s

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#FFF3E0;border-left:3px solid #E65100;border-radius:8px;margin:24px 0;">
      <tr>
        <td style="padding:16px 20px;">
          <div style="font-size:13px;color:#E65100;line-height:1.6;">
            <strong>⏱ This link expires in %d minutes.</strong><br>
            If you didn't request a password reset, you can safely ignore this email — your password will remain unchanged.
          </div>
        </td>
      </tr>
    </table>

    %s

    <p style="font-size:13px;color:#86868B;line-height:1.7;margin:0;">
      If the button doesn't work, copy and paste this URL into your browser:
    </p>
    <p style="font-size:12px;color:#1D1D1F;word-break:break-all;background-color:#FAFAFA;padding:12px;border-radius:8px;margin:12px 0 0;font-family:monospace;">
      %s
    </p>

    %s

    <p style="font-size:13px;color:#86868B;line-height:1.7;margin:0;">
      For your security, this password reset link can only be used once.
      If you need help, please contact us at
      <a href="mailto:zaemlifestyle@gmail.com" style="color:#1D1D1F;text-decoration:underline;">zaemlifestyle@gmail.com</a>.
    </p>
  `,
		Heading("Reset Your Password"),
		Paragraph(fmt.Sprintf("Hi %s,", safeName)),
		Button("Reset Password", resetURL),
		Divider(),
		resetURL,
	)

	preheader := fmt.Sprintf("Reset your ZAEM password — link valid for %d minutes", expiryMinutes)
	htmlBody = BaseLayout(subject, preheader, content)
	return
}