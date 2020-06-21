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
	fmt.Fprintf(w, "hello")

	var mentor models.Mentor
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mentor)
	if err != nil {
	  http.Error(w, err.Error(), 400)
	  return 		
	}

	
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		 http.Error(w, err.Error(), 500)
		 return
	}
	defer db.Close()

	err = db.Create(&models.Mentor{
				Name: Mentor.Name,
				Email: Mentor.Email,
				Github_handle: Mentor.Git_han,
				Access_token: Mentor.Acc_tok,
			}).Error;

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))

}

