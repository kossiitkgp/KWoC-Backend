package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONBody(r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(data)
	defer r.Body.Close()
	return err
}
