package routes

import (
	"fmt"
	"net/http"
)

func MentorReg(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
