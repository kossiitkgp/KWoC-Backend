package utils_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/kossiitkgp/kwoc-backend/v2/utils"
)

func TestJwtUtils(t *testing.T) {
	rand.Seed(time.Now().UnixMilli())

	testLoginFields := utils.LoginJwtFields{
		Username: fmt.Sprintf("testuser%d", rand.Int()),
	}

	// --- TEST SECRET KEY ERRORS ---
	_, err := utils.GenerateLoginJwtString(testLoginFields)

	if err != utils.ErrJwtSecretKeyNotFound {
		t.Errorf("Expected error %v. Got %v", utils.ErrJwtSecretKeyNotFound, err)
	}
	// --- TEST SECRET KEY ERRORS ---

	// Set a random test JWT key
	testJwtSecretKey := fmt.Sprintf("testkey%d", rand.Int())
	os.Setenv("JWT_SECRET_KEY", testJwtSecretKey)

	// --- TEST JWT GENERATION ---
	tokenString, err := utils.GenerateLoginJwtString(testLoginFields)

	if err != nil {
		t.Fatalf("Expected no error. Got %v.", err)
	}

	if tokenString == "" {
		t.Fatal("Expected a JWT string. Got empty string.")
	}
	// --- TEST JWT GENERATION ---

	// --- TEST JWT PARSING ---
	_, _, err = utils.ParseLoginJwtString("randominvalidtokenstring")

	if err == nil {
		t.Errorf("Expected error.")
	}

	token, claims, err := utils.ParseLoginJwtString(tokenString)

	if err != nil {
		t.Fatalf("Expected no error. Got %v.", err)
	}

	if !token.Valid {
		t.Fatalf("Expected valid token.")
	}

	if claims.Username != testLoginFields.Username {
		t.Errorf("Expected username `%s` in parsed claims. Got `%s`.", testLoginFields.Username, claims.Username)
	}
	// --- TEST JWT PARSING ---
}
