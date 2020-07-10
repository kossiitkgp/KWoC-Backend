package utils

import (
	"context"
	"net/http"
	"os"
	"fmt"

	logs "kwoc20-backend/utils/logs/pkg"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/go-kit/kit/log/level"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //For SQLite Dialect
)

//CtxUserString type for using with context
type CtxUserString string

//Claims jwt claim
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//LoginRequired Middleware to protect endpoints
func LoginRequired(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Bearer")
		if tokenStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jwtKey := []byte("OneRandomSecretKey!!@@!")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			level.Debug(logs.Logger).Log("message", "Invalid Token")
			return
		}

		ctx := r.Context()

		dbStr := "test.db"
		db, err := gorm.Open("sqlite3", dbStr)
		if err != nil {
			level.Error(logs.Logger).Log("error", "Failed to connect to the Database!")
			os.Exit(1)
		}
		defer db.Close()

		var user interface{} //Instance of Mentor/Mentee model
		newctx := context.WithValue(ctx, CtxUserString("user"), user)
		req := r.WithContext(newctx)

		next(w, req)
	}
}
