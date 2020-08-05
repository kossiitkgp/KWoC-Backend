package routes

import (
	"reflect"

	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"
)

func RegisterOAuth(r *mux.Router) {
	r.HandleFunc("", utils.JsonIO(controllers.UserOAuth, reflect.TypeOf(controllers.OAuthInput{}))).Methods("GET")

}
