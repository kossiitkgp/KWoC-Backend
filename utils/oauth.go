package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ErrClientIdNotFound = errors.New("ERROR: GITHUB OAUTH CLIENT ID NOT FOUND")
var ErrClientSecretNotFound = errors.New("ERROR: GITHUB OAUTH CLIENT SECRET NOT FOUND")
var ErrGithubAPIError = errors.New("ERROR: GITHUB API ERROR")

// Body fields for the request
type OAuthAccessReqFields struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

// Body fields for the response
type OAuthAccessResFields struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

func GetOauthAccessToken(code string) (string, error) {
	client_id := os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	client_secret := os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")

	if client_id == "" {
		return "", ErrClientIdNotFound
	}

	if client_secret == "" {
		return "", ErrClientSecretNotFound
	}

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make a request to the Github OAuth API to get the access token
	reqParams, err := json.Marshal(OAuthAccessReqFields{
		ClientId:     client_id,
		ClientSecret: client_secret,
		Code:         code,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		bytes.NewReader(reqParams),
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", nil
	}

	var resFields = OAuthAccessResFields{}

	err = json.NewDecoder(res.Body).Decode(&resFields)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	if resFields.Error != "" {
		return "", errors.New(resFields.Error)
	}

	return resFields.AccessToken, nil
}

type GHUserInfo struct {
	Username string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func GetOauthUserInfo(accessToken string) (*GHUserInfo, error) {
	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make a request to the Github OAuth API to get the username
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", accessToken),
	)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrGithubAPIError
	}

	var userInfo = GHUserInfo{}

	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
