package utils

import (
	"fmt"
	"log"
	"net/http"
)

func handleErr(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"message": "` + err.Error() + `"}`))

	fmt.Println("err test ", err.Error())

	log.Fatal("log is ", err)

}
