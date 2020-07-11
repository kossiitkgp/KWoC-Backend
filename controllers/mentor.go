package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"fmt"

	"kwoc20-backend/models"
	logs "kwoc20-backend/utils/logs/pkg"

	"github.com/jinzhu/gorm"
	"github.com/go-kit/kit/log/level"
)

// MentorOAuth Handler for Github OAuth of Mentor
func MentorOAuth(w http.ResponseWriter, r *http.Request) {
	// get the code from frontend
	var mentorOAuth1 interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error", fmt.Sprintf("%v",writeErr))
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
	err = json.Unmarshal(body, &mentorOAuth1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error", fmt.Sprintf("%v",writeErr))
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
	mentorOAuth, ok := mentorOAuth1.(map[string]interface{})
	if !ok {
		logErr := level.Error(logs.Logger).Log("error", "JSON mapping unsuccessful")
		if logErr != nil {
			panic("Log Error")
		}
		return
	}

	// using the code obtained from above to get AccessToken from Github
	req, err := json.Marshal(map[string]interface{}{
		"client_id":     os.Getenv("client_id"),
		"client_secret": os.Getenv("client_secret"),
		"code":          mentorOAuth["code"],
		"state":         os.Getenv("state"),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error", fmt.Sprintf("%v",writeErr))
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
	res, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer(req))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error", fmt.Sprintf("%v",writeErr))
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
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error", fmt.Sprintf("%v",writeErr))
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

	resBodyString := string(resBody)
	accessTokenPart := strings.Split(resBodyString, "&")[0]
	accessToken := strings.Split(accessTokenPart, "=")[1]

	// using the accessToken obtained above to get information about user
	client := &http.Client{}
	req1, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, writeErr := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if writeErr != nil {
			logErr := level.Warn(logs.Logger).Log("error", fmt.Sprintf("%v",writeErr))
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
	req1.Header.Add("Authorization", "token "+accessToken)
	res1, err := client.Do(req1)
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
	defer res1.Body.Close()
	resBody1, err := ioutil.ReadAll(res1.Body)
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

	var mentor1 interface{}
	err = json.Unmarshal(resBody1, &mentor1)
	if err != nil {
		logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",err))
			if logErr != nil {
				panic("Log Error")
			}
	}
	mentor, ok := mentor1.(map[string]interface{})
	if !ok {
		logErr := level.Error(logs.Logger).Log("error", "JSON mapping unsuccessful")
		if logErr != nil {
			panic("Log Error")
		}
		return
	}
	mentorData, err := json.Marshal(map[string]interface{}{
		"username": mentor["login"],
		"name":     mentor["name"],
		"email":    mentor["email"],
	})
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
	_, writeErr := w.Write(mentorData)
	if writeErr !=nil {
		logErr := level.Warn(logs.Logger).Log("error",fmt.Sprintf("%v",writeErr))
		if logErr != nil {
			panic("Log Error")
		}
	}

}

//MentorReg Handler for Registering Mentors
func MentorReg(w http.ResponseWriter, r *http.Request) {

	var mentor models.Mentor
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
	err = json.Unmarshal(body, &mentor)
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

	db.Create(&models.Mentor{
		Name:         mentor.Name,
		Email:        mentor.Email,
		GithubHandle: mentor.GithubHandle,
		AccessToken:  mentor.AccessToken,
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
}
