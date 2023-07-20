package utils

import (
	"encoding/json"
	"net/http"
)

type HTTPMessage struct {
	// HTTP Response Status Code
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func DecodeJSONBody(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	defer r.Body.Close()
	return err
}

func RespondWithJson(r *http.Request, w http.ResponseWriter, response any) {
	resJson, err := json.Marshal(response)

	if err != nil {
		LogErrAndRespond(r, w, err, "Error generating response JSON.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resJson)

	if err != nil {
		LogErr(r, err, "Error writing the response.")
		return
	}
}

func RespondWithHTTPMessage(r *http.Request, w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	res := HTTPMessage{
		Code:    status,
		Message: message,
	}
	RespondWithJson(r, w, res)
}
