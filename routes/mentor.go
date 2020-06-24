package routes

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	"kwoc20-backend/models"
)

func MentorReg(w http.ResponseWriter, r *http.Request) {

	var mentor models.Mentor
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mentor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
  		w.Write([]bytes(`{"message": "` + err.Error() + `"}`)
		return
	}

	
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
  		w.Write([]bytes(`{"message": "` + err.Error() + `"}`)
		return
	}
	defer db.Close()

	db.Create(&models.Mentor{
				Name: mentor.Name,
				Email: mentor.Email,
				Github_handle: mentor.Github_handle,
				Access_token: mentor.Access_token,
			})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
  	w.Write([]bytes(`{"message": "` + err.Error() + `"}`)
	return

}

