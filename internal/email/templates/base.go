package templates

import "fmt"

// BaseLayout wraps content with ZAEM's branded header + footer
// This is the design foundation for ALL emails
func BaseLayout(title, preheader, contentHTML string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="x-apple-disable-message-reformatting">
  <title>%s</title>
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings>
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <![endif]-->
</head>
<body style="margin:0;padding:0;background-color:#FAFAFA;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif;-webkit-font-smoothing:antialiased;">

  <!-- Preheader (hidden preview text) -->
  <div style="display:none;font-size:1px;color:#FAFAFA;line-height:1px;max-height:0;max-width:0;opacity:0;overflow:hidden;">
    %s
  </div>

  <!-- Outer table for email client compatibility -->
  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#FAFAFA;">
    <tr>
      <td align="center" style="padding:40px 16px;">

        <!-- Main container -->
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="max-width:600px;width:100%%;background-color:#FFFFFF;border-radius:16px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,0.04);">

          <!-- Header -->
          <tr>
            <td style="padding:40px 40px 24px;text-align:center;border-bottom:1px solid #F0F0F0;">
              <div style="font-family:Georgia,'Times New Roman',serif;font-size:28px;font-weight:400;letter-spacing:0.25em;color:#1D1D1F;margin:0;">
                ZAEM
              </div>
              <div style="font-size:9px;letter-spacing:0.4em;color:#86868B;margin-top:8px;text-transform:uppercase;">
                EST. 2026
              </div>
            </td>
          </tr>

          <!-- Content -->
          <tr>
            <td style="padding:40px;">
              %s
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="padding:32px 40px;background-color:#FAFAFA;border-top:1px solid #F0F0F0;text-align:center;">
              <div style="font-size:12px;color:#86868B;line-height:1.6;margin-bottom:16px;">
                Questions? Reply to this email or contact us at
                <a href="mailto:zaemlifestyle@gmail.com" style="color:#1D1D1F;text-decoration:underline;">zaemlifestyle@gmail.com</a>
              </div>

              <div style="margin:20px 0;">
                <a href="https://zaemstore.com" style="color:#86868B;text-decoration:none;font-size:11px;letter-spacing:0.1em;margin:0 12px;text-transform:uppercase;">Shop</a>
                <a href="https://zaemstore.com/account" style="color:#86868B;text-decoration:none;font-size:11px;letter-spacing:0.1em;margin:0 12px;text-transform:uppercase;">Account</a>
                <a href="https://zaemstore.com/contact" style="color:#86868B;text-decoration:none;font-size:11px;letter-spacing:0.1em;margin:0 12px;text-transform:uppercase;">Contact</a>
              </div>

              <div style="font-size:10px;color:#A1A1A6;margin-top:24px;line-height:1.6;">
                © 2026 ZAEM. All rights reserved.<br>
                Karachi, Pakistan
              </div>

              <div style="font-size:10px;color:#A1A1A6;margin-top:16px;">
                <a href="https://zaemstore.com/privacy" style="color:#A1A1A6;text-decoration:underline;">Privacy</a>
                &nbsp;·&nbsp;
                <a href="https://zaemstore.com/terms" style="color:#A1A1A6;text-decoration:underline;">Terms</a>
              </div>
            </td>
          </tr>

        </table>

        <!-- Below card spacing -->
        <div style="height:24px;"></div>

      </td>
    </tr>
  </table>

</body>
</html>`, title, preheader, contentHTML)
}

// Heading — reusable styled heading
func Heading(text string) string {
	return fmt.Sprintf(`<h1 style="font-family:Georgia,'Times New Roman',serif;font-size:28px;font-weight:400;color:#1D1D1F;margin:0 0 16px;line-height:1.3;">%s</h1>`, text)
}

// Paragraph — reusable styled paragraph
func Paragraph(text string) string {
	return fmt.Sprintf(`<p style="font-size:15px;color:#6E6E73;line-height:1.7;margin:0 0 20px;">%s</p>`, text)
}

// Button — branded CTA button
func Button(text, url string) string {
	return fmt.Sprintf(`<table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:32px 0;">
  <tr>
    <td style="background-color:#1D1D1F;border-radius:12px;">
      <a href="%s" style="display:inline-block;padding:16px 40px;color:#FFFFFF;text-decoration:none;font-size:12px;letter-spacing:0.2em;text-transform:uppercase;font-weight:500;">%s</a>
    </td>
  </tr>
</table>`, url, text)
}

// Divider — subtle line
func Divider() string {
	return `<div style="height:1px;background-color:#F0F0F0;margin:32px 0;"></div>`
}