package utils

import (
	"context"
	"net/http"
	"os"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
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
			http.Error(w, "Empty GET request", 400)
			LOG.Println("Empty Get request")
			return
		}

		jwtKey := []byte("OneRandomSecretKey!!@@!")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Empty GET request", 401)
			LOG.Println(err)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid Token", 401)
			LOG.Println("Invalid Token")
			return
		}

		ctx := r.Context()

		dbStr := "test.db"
		db, err := gorm.Open("sqlite3", dbStr)
		if err != nil {
			http.Error(w, "Failed to connect to the Database!", 500)
			LOG.Println(err)
			os.Exit(1)
		}
		defer db.Close()

		var user interface{} //Instance of Mentor/Mentee model
		newctx := context.WithValue(ctx, CtxUserString("user"), user)
		req := r.WithContext(newctx)

		next(w, req)
	}
}
