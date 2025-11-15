package controllers

import (
	"net/http"

	"github.com/kossiitkgp/kwoc-backend/v2/middleware"
	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

// IsAdmin godoc
//
//	@Summary		Checks if the logged in user is an admin
//	@Description	Returns true if the logged in user is an admin, false otherwise
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	map[string]bool	"is_admin status fetched successfully."
//	@Router			/isadmin [get]
func IsAdmin(w http.ResponseWriter, r *http.Request) {
	isAdmin := r.Context().Value(middleware.LOGIN_CTX_IS_ADMIN_KEY).(bool)
	utils.RespondWithJson(r, w, map[string]bool{"is_admin": isAdmin})
}
