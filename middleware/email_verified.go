package middleware

import (
	"net/http"
	"os"

	"github.com/kossiitkgp/kwoc-backend/v2/models"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

// WithVerifiedEmail blocks requests from users who have not verified their
// email address when EMAIL_VERIFICATION_REQUIRED is set to `true`.
//
// It must be placed inside WrapApp so that the database is available in the
// request context, and inside WithLogin so that the username is available.
func WithVerifiedEmail(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("EMAIL_VERIFICATION_REQUIRED") != "true" || !utils.EmailVerificationEnabled() {
			inner.ServeHTTP(w, r)
			return
		}

		app, ok := r.Context().Value(APP_CTX_KEY).(*App)
		if !ok || app == nil || app.Db == nil {
			utils.LogErrAndRespond(r, w, nil, "Error: Application context unavailable.", http.StatusInternalServerError)
			return
		}

		username, ok := r.Context().Value(LOGIN_CTX_USERNAME_KEY).(string)
		if !ok || username == "" {
			utils.LogWarnAndRespond(r, w, "Error: No logged in user found.", http.StatusUnauthorized)
			return
		}

		student := models.Student{}
		studentTx := app.Db.
			Select("id", "username", "email_verified").
			Where("username = ?", username).
			First(&student)

		if studentTx.Error == nil {
			if !student.EmailVerified {
				utils.LogWarnAndRespond(r, w, "Error: Email verification required.", http.StatusForbidden)
				return
			}

			inner.ServeHTTP(w, r)
			return
		}

		mentor := models.Mentor{}
		mentorTx := app.Db.
			Select("id", "username", "email_verified").
			Where("username = ?", username).
			First(&mentor)

		if mentorTx.Error == nil {
			if !mentor.EmailVerified {
				utils.LogWarnAndRespond(r, w, "Error: Email verification required.", http.StatusForbidden)
				return
			}

			inner.ServeHTTP(w, r)
			return
		}

		utils.LogWarnAndRespond(r, w, "Error: Registered student or mentor not found.", http.StatusForbidden)
	})
}
