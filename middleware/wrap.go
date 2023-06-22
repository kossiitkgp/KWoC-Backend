package middleware

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

// Reference: https://github.com/the-tikka-cafe/tyme/blob/main/backend/cmd/tyme/middleware.go
type App struct {
	Db *gorm.DB
}

type AppCtxKey string

const APP_CTX_KEY AppCtxKey = "app"

func WrapApp(app *App, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), APP_CTX_KEY, app)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
