package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONBody(r *http.Request, w http.ResponseWriter, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(data)
	defer r.Body.Close()
	if err != nil {
		LogErrAndRespond(r, w, err, "Error decoding JSON body.", http.StatusBadRequest)
		return err
	}
	return nil
}
