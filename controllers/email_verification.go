package controllers

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"

	"gorm.io/gorm"
)

type VerifyEmailReqFields struct {
	// Verification token from the email link.
	Token string `json:"token"`
	// Username to verify when using a code instead of a link.
	Username string `json:"username"`
	// Verification code from the email.
	Code string `json:"code"`
}

type VerifyEmailResFields struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	// `student` or `mentor`
	Role string `json:"role"`
	// Whether the account is now verified.
	EmailVerified bool `json:"email_verified"`
}

type ResendVerificationReqFields struct {
	// Optional new email address. If provided, the account email is updated
	// before the verification email is sent.
	Email string `json:"email"`
}

// VerifyEmail godoc
//
//	@Summary		Verify an email address
//	@Description	Verifies a student or mentor email address using the token from the verification link, or the username and code from the verification email.
//	@Accept			json
//	@Produce		json
//	@Param			request	body		VerifyEmailReqFields	true	"Verification token, or username and code"
//	@Success		200		{object}	VerifyEmailResFields	"Email verified successfully."
//	@Failure		400		{object}	utils.HTTPMessage		"Error decoding JSON body. | Invalid or expired verification code."
//	@Failure		500		{object}	utils.HTTPMessage		"Error verifying email address."
//	@Router			/verify/email/ [post]
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	reqFields := VerifyEmailReqFields{}
	if err := utils.DecodeJSONBody(r, &reqFields); err != nil {
		utils.LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
		return
	}

	verifiedUser, err := completeEmailVerification(r, db, reqFields.Token, reqFields.Username, reqFields.Code)
	if err != nil {
		status, message := verificationErrorResponse(err)
		if status == http.StatusInternalServerError {
			utils.LogErrAndRespond(r, w, err, "Error verifying email address.", status)
			return
		}

		utils.LogWarnAndRespond(r, w, message, status)
		return
	}

	utils.RespondWithJson(r, w, VerifyEmailResFields{
		Message:       "Email verified successfully.",
		Username:      verifiedUser.Username,
		Role:          verifiedUser.Role,
		EmailVerified: true,
	})
}

// VerifyEmailLink godoc
//
//	@Summary		Verify an email address from a link
//	@Description	Verifies a student or mentor email address using the token in the query string. This endpoint is used when the verification link in the email points directly to the backend.
//	@Produce		html
//	@Param			token	query		string	true	"Verification token"
//	@Success		200		{string}	string	"Email verified successfully."
//	@Failure		400		{string}	string	"Invalid or expired verification link."
//	@Failure		429		{string}	string	"Too many incorrect attempts."
//	@Failure		500		{string}	string	"Error verifying email address."
//	@Router			/verify/email/ [get]
func VerifyEmailLink(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		renderVerificationResultPage(
			w,
			http.StatusBadRequest,
			"Invalid verification link",
			"This verification link is missing its token. Please request a new verification email.",
		)
		return
	}

	verifiedUser, err := completeEmailVerification(r, db, token, "", "")
	if err != nil {
		status, message := verificationErrorResponse(err)
		if status == http.StatusInternalServerError {
			utils.LogErr(r, err, "Error verifying email address.")
		} else {
			utils.LogWarn(r, message)
		}

		renderVerificationResultPage(w, status, "Email verification failed", message)
		return
	}

	renderVerificationResultPage(
		w,
		http.StatusOK,
		"Email verified",
		fmt.Sprintf("Thanks, %s! Your %s account is now verified.", verifiedUser.Username, verifiedUser.Role),
	)
}

// ResendVerificationEmail godoc
//
//	@Summary		Resend the email verification code
//	@Description	Sends a new verification code and link to the logged in student or mentor. An optional new email address can be provided.
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ResendVerificationReqFields	false	"Optional new email address"
//	@Success		200		{object}	utils.HTTPMessage			"Verification email sent."
//	@Failure		400		{object}	utils.HTTPMessage			"Student or mentor does not exist. | Email address is already verified. | No email address is on record."
//	@Failure		500		{object}	utils.HTTPMessage			"Error sending verification email."
//	@Security		JWT
//	@Router			/verify/email/resend/ [post]
func ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
	db := app.Db

	loginUsername := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(string)

	reqFields := ResendVerificationReqFields{}
	// The body is optional, so only decode if one was provided.
	if r.Body != nil && r.ContentLength != 0 {
		if err := utils.DecodeJSONBody(r, &reqFields); err != nil && err != io.EOF {
			utils.LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
			return
		}
	}

	newEmail := strings.TrimSpace(reqFields.Email)

	student := models.Student{}
	studentTx := db.Where("username = ?", loginUsername).First(&student)
	if studentTx.Error != nil && studentTx.Error != gorm.ErrRecordNotFound {
		utils.LogErrAndRespond(r, w, studentTx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	if studentTx.Error == nil {
		if student.EmailVerified {
			utils.LogWarnAndRespond(r, w, "Email address is already verified.", http.StatusBadRequest)
			return
		}

		email := student.Email
		if newEmail != "" {
			email = newEmail
			if err := db.Model(&models.Student{}).
				Where("id = ?", student.ID).
				Update("email", newEmail).Error; err != nil {
				utils.LogErrAndRespond(r, w, err, "Error updating email address.", http.StatusInternalServerError)
				return
			}
		}

		if strings.TrimSpace(email) == "" {
			utils.LogWarnAndRespond(r, w, "No email address is on record. Add an email address first.", http.StatusBadRequest)
			return
		}

		sendVerificationEmailOrRespond(r, w, db, student.Username, student.Name, email, "student")
		return
	}

	mentor := models.Mentor{}
	mentorTx := db.Where("username = ?", loginUsername).First(&mentor)
	if mentorTx.Error == gorm.ErrRecordNotFound {
		utils.LogWarnAndRespond(r, w, "Student or mentor does not exist.", http.StatusBadRequest)
		return
	}
	if mentorTx.Error != nil {
		utils.LogErrAndRespond(r, w, mentorTx.Error, "Database error.", http.StatusInternalServerError)
		return
	}

	if mentor.EmailVerified {
		utils.LogWarnAndRespond(r, w, "Email address is already verified.", http.StatusBadRequest)
		return
	}

	email := mentor.Email
	if newEmail != "" {
		email = newEmail
		if err := db.Model(&models.Mentor{}).
			Where("id = ?", mentor.ID).
			Update("email", newEmail).Error; err != nil {
			utils.LogErrAndRespond(r, w, err, "Error updating email address.", http.StatusInternalServerError)
			return
		}
	}

	if strings.TrimSpace(email) == "" {
		utils.LogWarnAndRespond(r, w, "No email address is on record. Add an email address first.", http.StatusBadRequest)
		return
	}

	sendVerificationEmailOrRespond(r, w, db, mentor.Username, mentor.Name, email, "mentor")
}

// sendVerificationEmailOrRespond sends a verification email and writes the
// appropriate HTTP response.
func sendVerificationEmailOrRespond(
	r *http.Request,
	w http.ResponseWriter,
	db *gorm.DB,
	username string,
	name string,
	email string,
	role string,
) {
	if !utils.EmailVerificationEnabled() {
		utils.LogWarnAndRespond(r, w, "Email verification is disabled.", http.StatusBadRequest)
		return
	}

	if err := utils.SendAccountVerificationEmail(db, username, name, email, role); err != nil {
		if errors.Is(err, utils.ErrEmailNotConfigured) {
			utils.LogErrAndRespond(r, w, err, "Email sending is not configured.", http.StatusInternalServerError)
			return
		}

		utils.LogErrAndRespond(r, w, err, "Error sending verification email.", http.StatusInternalServerError)
		return
	}

	utils.RespondWithHTTPMessage(r, w, http.StatusOK, "Verification email sent. Please check your inbox.")
}

// sendVerificationEmailAfterRegistration sends a verification email during
// registration. Failures are logged, but do not fail the registration itself.
func sendVerificationEmailAfterRegistration(
	r *http.Request,
	db *gorm.DB,
	username string,
	name string,
	email string,
	role string,
) bool {
	if strings.TrimSpace(email) == "" {
		utils.LogWarn(r, "No email address provided; skipping verification email.")
		return false
	}

	if !utils.EmailVerificationEnabled() || !utils.IsEmailConfigured() {
		utils.LogWarn(r, "Email verification is disabled or email sending is not configured; skipping verification email.")
		return false
	}

	if err := utils.SendAccountVerificationEmail(db, username, name, email, role); err != nil {
		utils.LogErr(r, err, "Error sending verification email.")
		return false
	}

	return true
}

// sendRegistrationFailedEmail informs the user that registration failed. The
// failure is logged, but never changes the API response.
func sendRegistrationFailedEmail(
	r *http.Request,
	email string,
	name string,
	username string,
	role string,
	reason string,
) {
	if strings.TrimSpace(email) == "" || !utils.IsEmailConfigured() {
		return
	}

	if err := utils.SendRegistrationFailedEmail(email, name, username, role, reason); err != nil {
		utils.LogErr(r, err, "Error sending registration failed email.")
	}
}

// completeEmailVerification consumes a verification challenge, marks the
// account as verified, and sends the registration success email.
func completeEmailVerification(
	r *http.Request,
	db *gorm.DB,
	token string,
	username string,
	code string,
) (*utils.VerifiedUser, error) {
	verifiedUsername, err := utils.ConsumeEmailVerification(db, token, username, code)
	if err != nil {
		return nil, err
	}

	verifiedUser, err := utils.MarkUserEmailVerified(db, verifiedUsername)
	if err != nil {
		return nil, err
	}

	utils.LogInfo(r, "Email verified for `"+verifiedUser.Username+"`.")

	if verifiedUser.Email != "" {
		if err := utils.SendRegistrationSuccessEmail(
			verifiedUser.Email,
			verifiedUser.Name,
			verifiedUser.Username,
			verifiedUser.Role,
		); err != nil {
			utils.LogErr(r, err, "Error sending registration success email.")
		}
	}

	return verifiedUser, nil
}

// verificationErrorResponse maps a verification error to an HTTP status and a
// user-facing message.
func verificationErrorResponse(err error) (int, string) {
	switch {
	case errors.Is(err, utils.ErrVerificationMissingParams):
		return http.StatusBadRequest, "Provide a verification token, or a username and code."
	case errors.Is(err, utils.ErrVerificationCodeInvalid):
		return http.StatusBadRequest, "Invalid or expired verification code."
	case errors.Is(err, utils.ErrVerificationAttemptsExceeded):
		return http.StatusTooManyRequests, "Too many incorrect attempts. Request a new verification email."
	case errors.Is(err, utils.ErrVerificationExpired):
		return http.StatusBadRequest, "Verification link or code has expired. Request a new verification email."
	case errors.Is(err, utils.ErrVerificationAlreadyUsed):
		return http.StatusBadRequest, "This email address has already been verified."
	case errors.Is(err, utils.ErrVerificationNotFound):
		return http.StatusBadRequest, "Invalid or expired verification code."
	default:
		return http.StatusInternalServerError, "Error verifying email address."
	}
}

// renderVerificationResultPage writes a small HTML page for verification links
// that point directly to the backend.
func renderVerificationResultPage(w http.ResponseWriter, status int, heading string, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	dashboardURL := utils.FrontendURL()
	action := ""
	if dashboardURL != "" && status == http.StatusOK {
		action = fmt.Sprintf(
			`<p style="margin:24px 0 0 0;"><a href="%s" style="display:inline-block;padding:12px 22px;background:#111827;color:#ffffff;text-decoration:none;border-radius:8px;font-size:15px;font-weight:600;">Continue to KWoC</a></p>`,
			html.EscapeString(dashboardURL),
		)
	}

	_, _ = fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
</head>
<body style="margin:0;padding:0;background:#f5f5f5;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#1f2937;">
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background:#f5f5f5;padding:24px 12px;">
<tr>
<td align="center">
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="max-width:520px;background:#ffffff;border-radius:12px;border:1px solid #e5e7eb;">
<tr><td style="background:#111827;padding:20px 28px;color:#ffffff;font-size:20px;font-weight:700;">KWoC</td></tr>
<tr>
<td style="padding:32px 28px;">
<h1 style="margin:0 0 16px 0;font-size:24px;line-height:1.3;color:#111827;">%s</h1>
<p style="margin:0;font-size:16px;line-height:1.6;color:#4b5563;">%s</p>
%s
</td>
</tr>
</table>
</td>
</tr>
</table>
</body>
</html>`,
		html.EscapeString(heading),
		html.EscapeString(heading),
		html.EscapeString(message),
		action,
	)
}
