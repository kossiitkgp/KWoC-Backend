package routes

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	"kwoc20-backend/tests"
	"kwoc20-backend/utils"
)

func RegisterTest(r *mux.Router) {
	r.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		path := "./tests/oauth.html"
		http.ServeFile(w, r, path)
	})
	r.HandleFunc("/jsonio", utils.JsonIO(tests.JsonioTestFunc, reflect.TypeOf(tests.TestMarshalType{})))

}
