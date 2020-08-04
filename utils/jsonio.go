package utils

import (
    "io/ioutil"
    "encoding/json"
    "net/http"
    "reflect"
)

type ErrorMessage struct {
    Message string `json:"message"`
}

// JsonIO Middleware for JSON input and output
// Parameters of next: JSON as interface{}, Same request r (for other needs)
// Output of next: A struct pointer converted to interface{} and a bool ok
// Reference Usage:
// - Declare input and output structure as structs with json tags
// - Pass the input struct type as inputType
// - Use type switches to cast input interface{} to your Input struct
// - Cast response struct pointer to interface{}.
// See tests/jsonio.go for reference.
func JsonIO(next func(interface{}, *http.Request) (interface{}, bool), inputType reflect.Type) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        body, _ := ioutil.ReadAll(r.Body)

        jsonData := reflect.New(inputType)
        jsonPointer := jsonData.Interface()
        err := json.Unmarshal(body, jsonPointer)

        if err != nil {
            // Pass on silently
        }

        response, ok := next(jsonPointer, r)
        if !ok {
            w.WriteHeader(http.StatusBadRequest)
            w.Header().Set("Content-type", "application/json")
            w.Write([]byte(`{"message": "Invalid Request"}`))
            return
        }

        resBody, _ := json.Marshal(response)

        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-type", "application/json")
        w.Write(resBody)
    }
}

