package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// MentorOauth Handler for Github OAuth of Mentor
func UserOAuth(w http.ResponseWriter, r *http.Request) {
	// get the code from frontend
	var mentorOAuth1 interface{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mentorOAuth1)
	mentorOAuth, _ := mentorOAuth1.(map[string]interface{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// using the code obtained from above to get AccessToken from Github
	req, _ := json.Marshal(map[string]interface{}{
		"client_id":     os.Getenv("client_id"),
		"client_secret": os.Getenv("client_secret"),
		"code":          mentorOAuth["code"],
		"state":         mentorOAuth["state"],
	})
	res, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer(req))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)

	resBodyString := string(resBody)
	accessTokenPart := strings.Split(resBodyString, "&")[0]
	accessToken := strings.Split(accessTokenPart, "=")[1]

	// using the accessToken obtained above to get information about user
	client := &http.Client{}
	req1, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	req1.Header.Add("Authorization", "token "+accessToken)
	res1, err := client.Do(req1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer res1.Body.Close()
	resBody1, _ := ioutil.ReadAll(res1.Body)

	var mentor1 interface{}
	json.Unmarshal(resBody1, &mentor1)
	mentor, _ := mentor1.(map[string]interface{})
	mentorData, _ := json.Marshal(map[string]interface{}{
		"username": mentor["login"],
		"name":     mentor["name"],
		"email":    mentor["email"],
		"type":     mentorOAuth["state"],
	})

	w.WriteHeader(http.StatusOK)
	w.Write(mentorData)

}
