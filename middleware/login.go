package middleware

import (
	"context"
	"fmt"
	"kwoc-backend/utils"
	"net/http"
)

type LoginCtxKey string

var LOGIN_CTX_USERNAME_KEY string = "login_username"

// Session login middleware for incoming requests
func WithLogin(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Bearer")

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Error: No JWT session token found.")

			utils.LogWarn(r, "Unauthenticated request.")

			return
		}

		_, claims, err := utils.ParseLoginJwtString(tokenString)

		if err != nil {
			if err == utils.ErrJwtTokenInvalid {
				utils.LogWarn(r, "Invalid JWT Token Provided.")

				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Error: JWT session token invalid.")

				return
			}

			utils.LogErr(r, err, "Error parsing JWT string.")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error parsing JWT string.")

			return
		}

		reqContext := r.Context()
		newContext := context.WithValue(reqContext, LoginCtxKey(LOGIN_CTX_USERNAME_KEY), claims.LoginJwtFields.Username)

		inner.ServeHTTP(w, r.WithContext(newContext))
	})
}
