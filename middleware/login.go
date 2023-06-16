package middleware

import (
	"context"
	"fmt"
	"kwoc-backend/utils"
	"net/http"

	"github.com/rs/zerolog/log"
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

			log.Warn().Msgf(
				"%s %s %s",
				r.Method,
				r.RequestURI,
				"Unauthenticated request.",
			)

			return
		}

		_, claims, err := utils.ParseLoginJwtString(tokenString)

		if err != nil {
			if err == utils.ErrJwtTokenInvalid {
				log.Warn().Msgf(
					"%s %s %s",
					r.Method,
					r.RequestURI,
					"Invalid JWT Token Provided.",
				)

				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Error: JWT session token invalid.")

				return
			}

			log.Err(err).Msg("Error parsing JWT string.")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error parsing JWT string.")

			return
		}

		reqContext := r.Context()
		newContext := context.WithValue(reqContext, LoginCtxKey(LOGIN_CTX_USERNAME_KEY), claims.LoginJwtFields.Username)

		inner.ServeHTTP(w, r.WithContext(newContext))
	})
}
