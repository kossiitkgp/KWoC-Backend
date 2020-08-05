package tests

import (
    "kwoc20-backend/utils"
    "net/http"
)

type TestMarshalType struct {
    Message string `json:"message"`
}

func JsonioTestFunc(r interface{}, req *http.Request) (interface{}, bool) {
    switch r.(type) {
    case *TestMarshalType:
        break
    default:
        return &TestMarshalType{}, false
    }

    requestBody := r.(*TestMarshalType)
    utils.LOG.Println((*requestBody).Message)
    if (*requestBody).Message == "" {
        return &TestMarshalType{}, false
    }

    answer := &TestMarshalType{
        Message: "Message received",
    }

    return answer, true
}
