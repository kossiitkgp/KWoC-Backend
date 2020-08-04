package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"

	"kwoc20-backend/models"
	"kwoc20-backend/utils"
)

type OAuthInput struct {
	Code string `json:"code"`
	State string `json:"state"`
}

type OAuthOutput struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Email string `json:"email"`
	Type string `json:"type"`
	IsNewUser bool `json:"isNewUser"`
	JWT string `json:"jwt"`
	AccessToken string `json:"accessToken"`
}

// MentorOauth Handler for Github OAuth of Mentor
func UserOAuth(js interface{}, r *http.Request) (interface{}, bool) {
	// get the code from frontend
	mentorOAuth := &OAuthInput{
		Code: r.URL.Query().Get("code"),
		State: r.URL.Query().Get("state"),
	}

	if mentorOAuth.Code == "" || mentorOAuth.State == "" {
		return &utils.ErrorMessage{
			Message: "Type Mismatch",
		}, false
	}

	// using the code obtained from above to get AccessToken from Github
	req, _ := json.Marshal(map[string]interface{}{
		"client_id":     os.Getenv("client_id"),
		"client_secret": os.Getenv("client_secret"),
		"code":          mentorOAuth.Code,
		"state":         mentorOAuth.State,
	})
	res, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer(req))
	if err != nil {
		return &utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %s", err),
		}, false
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
		return &utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %+v", err),
		}, false
	}
	req1.Header.Add("Authorization", "token "+accessToken)
	res1, err := client.Do(req1)
	if err != nil {
		return &utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %+v", err),
		}, false
	}
	defer res1.Body.Close()

	resBody1, _ := ioutil.ReadAll(res1.Body)

	var userdata interface{}
	err = json.Unmarshal(resBody1, &userdata)
	if err != nil {
		return &utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %+v", err),
		}, false
	}

	user, _ := userdata.(map[string]interface{})

	gh_username, ok1 := user["login"].(string)
	gh_name, ok2 := user["name"].(string)
	gh_email, ok3 := user["email"].(string)

	if !ok1 {
		return &utils.ErrorMessage{
			Message: "GithubHandle not found",
		}, false
	}

	if !ok2 {
		gh_name = ""
	}

	if !ok3 {
		gh_email = ""
	}

	db, err := gorm.Open("sqlite3", "kwoc.db")
	if err != nil {
		return utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %+v", err),
		}, false
	}
	defer db.Close()

	chkUser := models.Mentor{}
	db.Where(&models.Mentor{GithubHandle: gh_username}).First(&chkUser)
	if chkUser.ID == 0 {
		// New User
		oauthdata := &OAuthOutput{
			Username: gh_username,
			Name: gh_name,
			Email: gh_username,
			Type: mentorOAuth.State,
			IsNewUser: true,
			JWT: "",
			AccessToken: accessToken,
		}
		utils.LOG.Println(fmt.Sprintf("New User: %+v", oauthdata))
		return oauthdata, true
	}

	// Returning user
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &utils.Claims{
		Username: chkUser.GithubHandle,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return &utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %+v", err),
		}, false
	}
	oauthdata := &OAuthOutput{
		Username: gh_username,
		Name: gh_name,
		Email: gh_email,
		Type: mentorOAuth.State,
		JWT: tokenStr,
		AccessToken: accessToken,
	}

	return oauthdata, true
}
