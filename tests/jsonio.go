package main

import (
    "kwoc20-backend/utils"
    "net/http"
    "reflect"
)

type MarshalType struct {
    Message string `json:"message"`
}

func testFunc(r interface{}, req *http.Request) (interface{}, bool) {
    switch r.(type) {
    case *MarshalType:
        break
    default:
        return &MarshalType{}, false
    }

    requestBody := r.(*MarshalType)
    utils.LOG.Println((*requestBody).Message)
    if (*requestBody).Message == "" {
        return &MarshalType{}, false
    }

    answer := &MarshalType{
        Message: "Message received",
    }

    return answer, true
}

func main() {
    http.HandleFunc("/", utils.JsonIO(testFunc, reflect.TypeOf(MarshalType{})))

    http.ListenAndServe(":2000", nil)
}
