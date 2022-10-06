package routes

import (
	"github.com/gorilla/mux"

	"kwoc20-backend/controllers"
	"kwoc20-backend/utils"
)

func RegisterOAuth(r *mux.Router) {
	r.HandleFunc("", utils.PermissiveCORSMiddleware(utils.JsonIO(controllers.UserOAuth))).Methods("POST")
}
