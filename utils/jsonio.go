package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

// JsonIO Middleware for JSON input and output
// Parameters of next: JSON as interface{}, Same request r (for other needs)
// Output of next: A struct pointer converted to interface{} and a int statusCode
// Reference Usage:
// - Declare input and output structure as structs with json tags
// - Pass the input struct type as inputType
// - Use type switches to cast input interface{} to your Input struct
// - Cast response struct pointer to interface{}.
// See tests/jsonio.go for reference.
func JsonIO(next func(map[string]interface{}, *http.Request) (interface{}, int)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recv := recover(); recv != nil {
				LogWarn(
					r,
					fmt.Sprintf("Kuch to locha hai\n%+v", recv),
				)
				debug.PrintStack()
				response := &ErrorMessage{
					Message: "Internal Server Error",
				}
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-type", "application/json")
				resBody, _ := json.Marshal(response)
				_, err := w.Write(resBody)
				if err != nil {
					LogErr(
						r,
						err,
						"ISSUE",
					)
				}
				return
			}
		}()

		body, _ := ioutil.ReadAll(r.Body)

		var jsonData1 interface{}
		_ = json.Unmarshal(body, &jsonData1)
		var jsonData map[string]interface{}
		if jsonData1 == nil {
			jsonData = make(map[string]interface{})
		} else {
			jsonData = jsonData1.(map[string]interface{})
		}

		response, statusCode := next(jsonData, r)
		// if statusCode is not in 200s, in case of error
		if statusCode/100 > 2 {
			LogWarn(
				r,
				fmt.Sprintf("%+v", response),
			)

			w.WriteHeader(statusCode)
			w.Header().Set("Content-type", "application/json")

			_, err := w.Write([]byte(`{"message": "Invalid Request"}`))
			if err != nil {
				LogErr(
					r,
					err,
					"ISSUE",
				)
			}
			return
		}

		resBody, _ := json.Marshal(response)

		w.WriteHeader(statusCode)
		w.Header().Set("Content-type", "application/json")
		_, _ = w.Write(resBody)
	}
}
