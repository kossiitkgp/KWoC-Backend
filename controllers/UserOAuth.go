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

	"github.com/golang-jwt/jwt/v4"

	"kwoc20-backend/models"
	"kwoc20-backend/utils"
)

// Handler for UserOAuth
func UserOAuth(js map[string]interface{}, r *http.Request) (interface{}, int) {

	// return error if no state or no code
	if js["code"] == "" || js["state"] == "" {
		return "type mismatch", 400
	}

	// using the code obtained from above to get AccessToken from Github
	req, _ := json.Marshal(map[string]interface{}{
		"client_id":     os.Getenv("client_id"),
		"client_secret": os.Getenv("client_secret"),
		"code":          js["code"],
		"state":         js["state"],
	})
	res, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer(req))
	if err != nil {
		return fmt.Sprintf("Error occurred: %s", err), 500
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
		return fmt.Sprintf("Error occurred: %+v", err), 500
	}
	req1.Header.Add("Authorization", "token "+accessToken)
	res1, err := client.Do(req1)
	if err != nil {
		return fmt.Sprintf("Error occurred: %+v", err), 500
	}
	defer res1.Body.Close()

	resBody1, _ := ioutil.ReadAll(res1.Body)

	var userdata interface{}
	err = json.Unmarshal(resBody1, &userdata)
	if err != nil {
		return &utils.ErrorMessage{
			Message: fmt.Sprintf("Error occurred: %+v", err),
		}, 500
	}

	user, _ := userdata.(map[string]interface{})

	gh_username, ok1 := user["login"].(string)
	gh_name, ok2 := user["name"].(string)
	gh_email, ok3 := user["email"].(string)

	utils.LogInfo(
		r,
		fmt.Sprintf("%+v %+v %+v\n", ok1, ok2, ok3),
	)

	if !ok1 {
		return &utils.ErrorMessage{
			Message: "GithubHandle not found",
		}, 500
	}

	if !ok2 {
		gh_name = ""
	}

	if !ok3 {
		gh_email = ""
	}

	db := utils.GetDB()
	defer db.Close()

	var isNewUser uint
	if js["state"] == "mentor" {
		chkUser := models.Mentor{}
		db.Where(&models.Mentor{Username: gh_username}).First(&chkUser)
		isNewUser = chkUser.ID
	} else {
		chkUser := models.Student{}
		db.Where(&models.Student{Username: gh_username}).First(&chkUser)
		isNewUser = chkUser.ID
	}

	// Creating a JWT token
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	expirationTime := time.Now().Add(50 * 24 * time.Hour)
	claims := &utils.Claims{
		Username: gh_username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return fmt.Sprintf("Error occurred: %+v", err), 500
	}

	if isNewUser == 0 {
		// New User
		resNewUser := map[string]interface{}{
			"username":    gh_username,
			"name":        gh_name,
			"email":       gh_email,
			"type":        js["state"],
			"isNewUser":   1,
			"jwt":         tokenStr,
			"accessToken": accessToken,
		}

		utils.LogInfo(
			r,
			fmt.Sprintf(
				"New User: %+v",
				resNewUser,
			),
		)
		return resNewUser, 200
	}

	resOldUser := map[string]interface{}{
		"username":    gh_username,
		"name":        gh_name,
		"email":       gh_email,
		"type":        js["state"],
		"isNewUser":   0,
		"jwt":         tokenStr,
		"accessToken": accessToken,
	}
	return resOldUser, 200

}
