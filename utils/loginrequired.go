package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	// _ "github.com/jinzhu/gorm/dialects/mysql" // For MySQL Dialect
)

// CtxUserString type for using with context
type CtxUserString string

// Claims jwt claim
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// LoginRequired Middleware to protect endpoints
func LoginRequired(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// If development, bypass the bearer token with the developer's github username
		if os.Getenv("MODE") == "dev" {
			ctx := r.Context()
			newctx := context.WithValue(ctx, CtxUserString("user"), os.Getenv("DEV_USERNAME"))
			req := r.WithContext(newctx)

			next(w, req)
			return
		}

		bearer_token := r.Header.Get("Authorization")
		tokenStrings := strings.Split(bearer_token, " ")
		fmt.Print(tokenStrings)
		if len(tokenStrings) != 2 {
			http.Error(w, "Invalid Authorization Token", http.StatusUnauthorized)
			return
		}
		tokenStr := tokenStrings[1]

		if tokenStr == "" {
			http.Error(w, "Empty GET request", http.StatusUnauthorized)
			LOG.Println("Empty Get request")
			return
		}

		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Empty GET request", http.StatusUnauthorized)
			LOG.Println(err)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			LOG.Println("Invalid Token")
			return
		}

		ctx := r.Context()
		newctx := context.WithValue(ctx, CtxUserString("user"), claims.Username)
		req := r.WithContext(newctx)
		next(w, req)
	}
}
