package middleware

import (
	"context"
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

type LoginCtxKey string

const LOGIN_CTX_USERNAME_KEY LoginCtxKey = "login_username"

// Session login middleware for incoming requests
func WithLogin(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Bearer")

		if tokenString == "" {
			utils.LogWarnAndRespond(r, w, "Error: No JWT session token found.", http.StatusUnauthorized)
			return
		}

		_, claims, err := utils.ParseLoginJwtString(tokenString)

		if err != nil {
			if err == utils.ErrJwtTokenInvalid {
				utils.LogErrAndRespond(r, w, err, "Error: JWT session token invalid.", http.StatusUnauthorized)
				return
			}

			if err == utils.ErrJwtTokenExpired {
				utils.LogErrAndRespond(r, w, err, "Error: JWT session token expired.", http.StatusUnauthorized)
				return
			}

			utils.LogErrAndRespond(r, w, err, "Error parsing JWT string.", http.StatusInternalServerError)
			return
		}

		reqContext := r.Context()
		newContext := context.WithValue(reqContext, LoginCtxKey(LOGIN_CTX_USERNAME_KEY), claims.LoginJwtFields.Username)

		inner.ServeHTTP(w, r.WithContext(newContext))
	})
}
