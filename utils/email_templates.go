package utils

import (
	"fmt"
	"html"
	"strings"
)

func roleDisplayName(role string) string {
	switch strings.ToLower(role) {
	case "mentor":
		return "Mentor"
	case "admin":
		return "Admin"
	default:
		return "Student"
	}
}

func greeting(name string, username string) string {
	if strings.TrimSpace(name) != "" {
		return name
	}

	if strings.TrimSpace(username) != "" {
		return username
	}

	return "there"
}

// SendVerificationEmail sends a verification code and link to the user.
func SendVerificationEmail(to string, name string, username string, code string, link string, role string) error {
	expiry := verificationExpiry()
	expiryMinutes := int(expiry.Minutes())
	if expiryMinutes <= 0 {
		expiryMinutes = defaultVerificationExpiryMinutes
	}

	subject := "Verify your email for KWoC"
	text := fmt.Sprintf(
		"Hi %s,\n\nUse the code below to verify the email address for your KWoC %s account (%s):\n\n%s\n\nThis code expires in %d minutes.\n\nYou can also verify by opening this link:\n%s\n\nIf you did not create a KWoC account, you can ignore this email.\n\nNeed help? Contact %s.\n",
		greeting(name, username),
		strings.ToLower(roleDisplayName(role)),
		username,
		code,
		expiryMinutes,
		link,
		SupportEmail(),
	)

	body := fmt.Sprintf(`
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Hi %s,</p>
<p style="margin:0 0 20px 0;font-size:16px;line-height:1.6;">Use the code below to verify the email address for your KWoC %s account <strong>%s</strong>.</p>
<div style="margin:0 0 20px 0;padding:18px 24px;background:#f3f4f6;border-radius:10px;text-align:center;">
<span style="font-family:'SFMono-Regular',Consolas,'Liberation Mono',Menlo,monospace;font-size:32px;font-weight:700;letter-spacing:10px;color:#111827;">%s</span>
</div>
<p style="margin:0 0 16px 0;font-size:14px;line-height:1.6;color:#4b5563;">This code expires in <strong>%d minutes</strong>.</p>
%s
<p style="margin:20px 0 0 0;font-size:14px;line-height:1.6;color:#4b5563;">If you did not create a KWoC account, you can safely ignore this email.</p>
`,
		html.EscapeString(greeting(name, username)),
		html.EscapeString(strings.ToLower(roleDisplayName(role))),
		html.EscapeString(username),
		html.EscapeString(code),
		expiryMinutes,
		emailButton(link, "Verify email address"),
	)

	return SendEmail(Mail{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    renderEmailTemplate(subject, "Your KWoC verification code is inside.", "Verify your email", body),
	})
}

// SendRegistrationSuccessEmail notifies a user that their registration and
// email verification are complete.
func SendRegistrationSuccessEmail(to string, name string, username string, role string) error {
	roleName := roleDisplayName(role)
	subject := fmt.Sprintf("Welcome to KWoC, %s!", greeting(name, username))
	dashboardURL := FrontendURL()
	if dashboardURL != "" {
		dashboardURL += "/" + strings.ToLower(roleName) + "/dashboard"
	}

	text := fmt.Sprintf(
		"Hi %s,\n\nYour KWoC %s registration is complete. Your email address has been verified.\n\nUsername: %s\n\nYou can now sign in and explore KWoC.\n%s\n\nNeed help? Contact %s.\n",
		greeting(name, username),
		strings.ToLower(roleName),
		username,
		dashboardURL,
		SupportEmail(),
	)

	body := fmt.Sprintf(`
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Hi %s,</p>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Your KWoC %s registration is complete and your email address has been verified.</p>
<div style="margin:0 0 20px 0;padding:16px 20px;background:#f0fdf4;border:1px solid #bbf7d0;border-radius:10px;font-size:15px;line-height:1.6;color:#166534;">
<strong>Username:</strong> %s<br>
<strong>Status:</strong> Email verified
</div>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">You can now sign in and get started.</p>
%s
`,
		html.EscapeString(greeting(name, username)),
		html.EscapeString(strings.ToLower(roleName)),
		html.EscapeString(username),
		emailButton(dashboardURL, "Go to dashboard"),
	)

	return SendEmail(Mail{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    renderEmailTemplate(subject, "Your KWoC registration is complete.", "Registration successful", body),
	})
}

// SendRegistrationFailedEmail notifies a user that their registration could
// not be completed, along with the reason and a way to get help.
func SendRegistrationFailedEmail(to string, name string, username string, role string, reason string) error {
	roleName := roleDisplayName(role)
	subject := "Your KWoC registration could not be completed"

	text := fmt.Sprintf(
		"Hi %s,\n\nWe could not complete the registration for the KWoC %s account `%s`.\n\nReason: %s\n\nIf you think this is a mistake, or if you need help, contact %s.\n",
		greeting(name, username),
		strings.ToLower(roleName),
		username,
		reason,
		SupportEmail(),
	)

	body := fmt.Sprintf(`
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Hi %s,</p>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">We could not complete the registration for the KWoC %s account <strong>%s</strong>.</p>
<div style="margin:0 0 20px 0;padding:16px 20px;background:#fef2f2;border:1px solid #fecaca;border-radius:10px;font-size:15px;line-height:1.6;color:#991b1b;">
<strong>Reason:</strong> %s
</div>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">If you think this is a mistake, or if you need help resolving it, please contact us.</p>
%s
`,
		html.EscapeString(greeting(name, username)),
		html.EscapeString(strings.ToLower(roleName)),
		html.EscapeString(username),
		html.EscapeString(reason),
		emailButton("mailto:"+SupportEmail(), "Contact support"),
	)

	return SendEmail(Mail{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    renderEmailTemplate(subject, "Action may be needed on your KWoC registration.", "Registration not completed", body),
	})
}

// SendProjectApprovedEmail tells a mentor that their project was approved.
func SendProjectApprovedEmail(to string, name string, username string, projectName string) error {
	subject := fmt.Sprintf("Your KWoC project \"%s\" has been approved", projectName)
	dashboardURL := FrontendURL()
	if dashboardURL != "" {
		dashboardURL += "/mentor/dashboard"
	}

	text := fmt.Sprintf(
		"Hi %s,\n\nGood news! Your project \"%s\" has been approved and is now live on KWoC.\n\nStudents can now find and contribute to your project.\n%s\n\nNeed help? Contact %s.\n",
		greeting(name, username),
		projectName,
		dashboardURL,
		SupportEmail(),
	)

	body := fmt.Sprintf(`
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Hi %s,</p>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Good news! Your project <strong>%s</strong> has been approved and is now live on KWoC.</p>
<div style="margin:0 0 20px 0;padding:16px 20px;background:#f0fdf4;border:1px solid #bbf7d0;border-radius:10px;font-size:15px;line-height:1.6;color:#166534;">
<strong>Status:</strong> Approved<br>
<strong>Project:</strong> %s
</div>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Students can now find and contribute to your project.</p>
%s
`,
		html.EscapeString(greeting(name, username)),
		html.EscapeString(projectName),
		html.EscapeString(projectName),
		emailButton(dashboardURL, "View project"),
	)

	return SendEmail(Mail{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    renderEmailTemplate(subject, "Your KWoC project has been approved.", "Project approved", body),
	})
}

// SendActionRequiredEmail sends an email telling the user that they need to do
// something for their account or project.
func SendActionRequiredEmail(to string, name string, username string, heading string, message string, actionURL string, actionLabel string) error {
	subject := heading
	if !strings.HasPrefix(strings.ToLower(subject), "action required") {
		subject = "Action required: " + subject
	}

	text := fmt.Sprintf(
		"Hi %s,\n\n%s\n\n%s\n\nNeed help? Contact %s.\n",
		greeting(name, username),
		heading,
		message,
		SupportEmail(),
	)

	body := fmt.Sprintf(`
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">Hi %s,</p>
<div style="margin:0 0 20px 0;padding:16px 20px;background:#fffbeb;border:1px solid #fde68a;border-radius:10px;font-size:15px;line-height:1.6;color:#92400e;">
%s
</div>
<p style="margin:0 0 16px 0;font-size:16px;line-height:1.6;">%s</p>
%s
`,
		html.EscapeString(greeting(name, username)),
		html.EscapeString(heading),
		html.EscapeString(message),
		emailButton(actionURL, actionLabel),
	)

	return SendEmail(Mail{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    renderEmailTemplate(subject, heading, heading, body),
	})
}

func renderEmailTemplate(title string, preheader string, heading string, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
</head>
<body style="margin:0;padding:0;background:#f5f5f5;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#1f2937;">
<div style="display:none;max-height:0;overflow:hidden;opacity:0;">%s</div>
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background:#f5f5f5;padding:24px 12px;">
<tr>
<td align="center">
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="max-width:560px;background:#ffffff;border-radius:12px;overflow:hidden;border:1px solid #e5e7eb;">
<tr>
<td style="background:#111827;padding:20px 28px;color:#ffffff;font-size:20px;font-weight:700;letter-spacing:0.3px;">KWoC</td>
</tr>
<tr>
<td style="padding:32px 28px 8px 28px;">
<h1 style="margin:0 0 18px 0;font-size:24px;line-height:1.3;color:#111827;">%s</h1>
%s
</td>
</tr>
<tr>
<td style="padding:24px 28px 32px 28px;color:#6b7280;font-size:13px;line-height:1.6;border-top:1px solid #f3f4f6;">
<p style="margin:0;">Need help? Contact <a href="mailto:%s" style="color:#2563eb;">%s</a>.</p>
<p style="margin:8px 0 0 0;">You are receiving this email because of your KWoC account. If this wasn't you, you can ignore this message.</p>
</td>
</tr>
</table>
</td>
</tr>
</table>
</body>
</html>`,
		html.EscapeString(title),
		html.EscapeString(preheader),
		html.EscapeString(heading),
		body,
		html.EscapeString(SupportEmail()),
		html.EscapeString(SupportEmail()),
	)
}

func emailButton(url string, label string) string {
	if strings.TrimSpace(url) == "" {
		return ""
	}

	return fmt.Sprintf(
		`<p style="margin:24px 0;"><a href="%s" style="display:inline-block;padding:12px 22px;background:#111827;color:#ffffff;text-decoration:none;border-radius:8px;font-size:15px;font-weight:600;">%s</a></p>`,
		html.EscapeString(url),
		html.EscapeString(label),
	)
}
