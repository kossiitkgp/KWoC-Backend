package tests

import "fmt"

func main() {
	fmt.Println("Testing")
}

// import (
//     "kwoc20-backend/utils"
//     "net/http"
// )

// type TestMarshalType struct {
//     Message string `json:"message"`
// }

// func JsonioTestFunc(r interface{}, req *http.Request) (interface{}, int) {
//     switch r.(type) {
//     case *TestMarshalType:
//         break
//     default:
//         return &TestMarshalType{}, 500
//     }

//     requestBody := r.(*TestMarshalType)
//     utils.LOG.Println((*requestBody).Message)
//     if (*requestBody).Message == "" {
//         return &TestMarshalType{}, 400
//     }

//     answer := &TestMarshalType{
//         Message: "Message received",
//     }

//     return answer, 200
// }
