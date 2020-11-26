package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // For MySQL Dialect
)

//CtxUserString type for using with context
type CtxUserString string

//Claims jwt claim
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// generates a random string with length given as a parameter
func generateRandomString(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyz"
	// characters to be used in the random string can be added in this string
	n := len(characters)

	rand.Seed(time.Now().UnixNano())
	random_string := "" 
	for i:=1; i<=length; i++ {
		random_number := rand.Intn(n)
		random_string += string(characters[int(random_number)])
	}
	return random_string
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
		secret_string := generateRandomString(6)
		jwtKey := []byte(secret_string)

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
		DatabaseUsername := os.Getenv("DATABASE_USERNAME")
		DatabasePassword := os.Getenv("DATABASE_PASSWORD")
		DatabaseName := os.Getenv("DATABASE_NAME")
		DatabaseHost := os.Getenv("DATABASE_HOST")
		DatabasePort := os.Getenv("DATABASE_PORT")

		DatabaseURI := fmt.Sprintf(
			"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			DatabaseUsername,
			DatabasePassword,
			DatabaseHost,
			DatabasePort,
			DatabaseName,
		)

		db, err := gorm.Open("mysql", DatabaseURI)
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
