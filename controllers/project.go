package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/go-kit/kit/log/level"
	logs "kwoc20-backend/utils/logs/pkg"

	"kwoc20-backend/models"
)

//ProjectReg endpoint to register project details
func ProjectReg(w http.ResponseWriter, r *http.Request) {

	var project models.Project
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}
		return
	}
	err = json.Unmarshal(body, &project)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}
		return
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}
		return
	}
	defer db.Close()

	err = db.Create(&models.Project{
		Name:       project.Name,
		Desc:       project.Desc,
		Tags:       project.Tags,
		RepoLink:   project.RepoLink,
		ComChannel: project.ComChannel,
	}).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write([]byte(`{"message": "success"}`))
	if writeErr !=nil {
		logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
		if logErr != nil {
			panic("Log Error")
		}
	}
	logErr := level.Info(logs.Logger).Log("message", "Succesfully registered project details")
	if logErr != nil {
		panic("Log Error")
	}
}

//ProjectGet endpoint to fetch all projects
// INCOMPLETE BECAUSE MENTOR STILL NEEDS TO BE ADDED
func ProjectGet(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}

		return
	}
	defer db.Close()

	var projects []models.Project
	err = db.Find(&projects).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(projects)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr !=nil {
			logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
			if logErr != nil {
				panic("Log Error")
			}
		}
		logErr := level.Error(logs.Logger).Log("error", fmt.Sprintf("%v",err))
		if logErr != nil {
			panic("Log Error")
		}
		return
	}

}
