package middleware

import (
	"context"
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

type LoginCtxKey string
type AdminCtxKey string

const LOGIN_CTX_USERNAME_KEY LoginCtxKey = "login_username"
const LOGIN_CTX_IS_ADMIN_KEY AdminCtxKey = "login_is_admin"

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

		username := claims.LoginJwtFields.Username

		isAdmin, adminErr := utils.IsUserAdmin(username)
		if adminErr != nil {
			utils.LogErrAndRespond(r, w, adminErr, "Error checking if user is admin.", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, LOGIN_CTX_USERNAME_KEY, username)
		ctx = context.WithValue(ctx, LOGIN_CTX_IS_ADMIN_KEY, isAdmin)

		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
