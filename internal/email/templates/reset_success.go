package templates

import (
	"fmt"
	stdhtml "html"
)

// PasswordResetSuccessEmail — sent after password successfully changed
func PasswordResetSuccessEmail(userName, loginURL string) (subject, htmlBody string) {
	subject = "Your ZAEM Password Has Been Changed"

	safeName := stdhtml.EscapeString(userName)

	content := fmt.Sprintf(`
    %s

    %s

    <p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">
      Your ZAEM account password was successfully changed. You can now use your new password to sign in.
    </p>

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#E8F5E9;border-left:3px solid #2E7D32;border-radius:8px;margin:24px 0;">
      <tr>
        <td style="padding:16px 20px;">
          <div style="font-size:13px;color:#2E7D32;line-height:1.6;">
            <strong>✓ Your password was updated successfully.</strong>
          </div>
        </td>
      </tr>
    </table>

    %s

    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#FFEBEE;border-left:3px solid #C62828;border-radius:8px;margin:24px 0;">
      <tr>
        <td style="padding:16px 20px;">
          <div style="font-size:13px;color:#C62828;line-height:1.6;">
            <strong>⚠ Didn't make this change?</strong><br>
            Contact us immediately at <a href="mailto:zaemlifestyle@gmail.com" style="color:#C62828;text-decoration:underline;">zaemlifestyle@gmail.com</a>.
          </div>
        </td>
      </tr>
    </table>
  `,
		Heading("Password Changed Successfully"),
		Paragraph(fmt.Sprintf("Hi %s,", safeName)),
		Button("Sign In", loginURL),
	)

	preheader := "Your ZAEM password has been changed"
	htmlBody = BaseLayout(subject, preheader, content)
	return
}