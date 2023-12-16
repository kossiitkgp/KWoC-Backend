package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var ErrJwtSecretKeyNotFound = errors.New("ERROR: JWT SECRET KEY NOT FOUND")
var ErrJwtTokenExpired = errors.New("ERROR: JWT TOKEN EXPIRED")
var ErrJwtTokenInvalid = errors.New("ERROR: JWT TOKEN INVALID")

func getJwtKey() (string, error) {
	jwtKey := os.Getenv("JWT_SECRET_KEY")

	if jwtKey == "" {
		return "", ErrJwtSecretKeyNotFound
	}

	return jwtKey, nil
}

func jwtKeyFunc(*jwt.Token) (interface{}, error) {
	key, err := getJwtKey()

	if err != nil {
		return nil, err
	}

	return []byte(key), err
}

type LoginJwtFields struct {
	Username string `json:"username"`
}

type LoginJwtClaims struct {
	LoginJwtFields
	jwt.RegisteredClaims
}

func ParseLoginJwtString(tokenString string) (*jwt.Token, *LoginJwtClaims, error) {
	var loginClaims = LoginJwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &loginClaims, jwtKeyFunc)

	if err.Error() == fmt.Sprintf("%s: %s", jwt.ErrTokenInvalidClaims.Error(), jwt.ErrTokenExpired.Error()) {
		return nil, nil, ErrJwtTokenExpired
	}

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, ErrJwtTokenInvalid
	}

	return token, &loginClaims, nil
}

func GenerateLoginJwtString(loginJwtFields LoginJwtFields) (string, error) {
	issue_time := time.Now()

	// Get the amount of time for which the generated JWT will be valid
	jwtValidityTimeEnvVar := os.Getenv("JWT_VALIDITY_TIME")
	jwtValidityTime, err := strconv.Atoi(jwtValidityTimeEnvVar)

	if err != nil {
		// Default of 30 days
		jwtValidityTime = 24 * 30

		log.Warn().Msgf("Could not parse JWT validity time from the environment. Set to default of %d hours.", jwtValidityTime)
	}

	claims := &LoginJwtClaims{
		LoginJwtFields: loginJwtFields,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issue_time),
			NotBefore: jwt.NewNumericDate(issue_time),
			// Valid for 30 days
			ExpiresAt: jwt.NewNumericDate(issue_time.Add(time.Duration(jwtValidityTime) * time.Hour)),
		},
	}

	signingKey, err := getJwtKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
