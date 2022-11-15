package utils

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
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
		tokenStr := r.Header.Get("Bearer")
		if tokenStr == "" {
			http.Error(w, "Empty GET request", 400)

			LogWarn(
				r,
				"Empty GET Request",
			)
			return
		}

		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Empty GET request", http.StatusUnauthorized)

			LogWarn(
				r,
				"Empty GET Request",
			)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)

			LogWarn(
				r,
				"Invalid Token",
			)
			return
		}

		ctx := r.Context()
		// DatabaseUsername := os.Getenv("DATABASE_USERNAME")
		// DatabasePassword := os.Getenv("DATABASE_PASSWORD")
		// DatabaseName := os.Getenv("DATABASE_NAME")
		// DatabaseHost := os.Getenv("DATABASE_HOST")
		// DatabasePort := os.Getenv("DATABASE_PORT")

		// DatabaseURI := fmt.Sprintf(
		// 	"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		// 	DatabaseUsername,
		// 	DatabasePassword,
		// 	DatabaseHost,
		// 	DatabasePort,
		// 	DatabaseName,
		// )

		// db, err := gorm.Open("mysql", DatabaseURI)
		// if err != nil {
		// 	http.Error(w, "Failed to connect to the Database!", 500)
		// 	LOG.Println(err)
		// 	os.Exit(1)
		// }
		// defer db.Close()

		// var user interface{} //Instance of Mentor/Mentee model
		newctx := context.WithValue(ctx, CtxUserString("user"), claims.Username)
		req := r.WithContext(newctx)

		next(w, req)
	}
}
