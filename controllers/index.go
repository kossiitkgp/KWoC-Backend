// Handler function for routes
package controllers

import (
	"fmt"
	"net/http"
)

// Index godoc
//
//	@Summary		Checks api health
//	@Description	Returns a simple string to check if the api is up and running
//	@Accept			plain
//	@Produce		plain
//	@Success		200	{string}	string	"Hello from KOSS Backend!"
//	@Router			/api [get]
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from KOSS Backend!")
}
