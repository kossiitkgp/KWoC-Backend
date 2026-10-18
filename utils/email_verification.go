package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kossiitkgp/kwoc-backend/v2/models"

	"gorm.io/gorm"
)

const (
	defaultVerificationExpiryMinutes = 30
	defaultVerificationMaxAttempts   = 5
	verificationCodeDigits           = 6
	verificationTokenBytes           = 32
)

var (
	ErrVerificationNotFound         = errors.New("verification request not found")
	ErrVerificationExpired          = errors.New("verification request expired")
	ErrVerificationAlreadyUsed      = errors.New("verification request already used")
	ErrVerificationAttemptsExceeded = errors.New("too many verification attempts")
	ErrVerificationCodeInvalid      = errors.New("invalid verification code")
	ErrVerificationMissingParams    = errors.New("verification token or username and code are required")
)

// VerificationChallenge contains the plaintext values that are emailed to the
// user. Only their hashes are persisted.
type VerificationChallenge struct {
	Record    *models.EmailVerification
	Code      string
	Token     string
	Link      string
	ExpiresIn time.Duration
}

// VerifiedUser is the account that was marked as verified.
type VerifiedUser struct {
	Username string
	Name     string
	Email    string
	Role     string
}

func VerificationMaxAttempts() int {
	return envIntOrDefault("EMAIL_VERIFICATION_MAX_ATTEMPTS", defaultVerificationMaxAttempts)
}

// EmailVerificationEnabled reports whether verification emails should be sent.
func EmailVerificationEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("EMAIL_VERIFICATION_ENABLED")))
	return value == "" || value == "true" || value == "1" || value == "yes"
}

// verificationExpiry returns the lifetime of a verification challenge.
func verificationExpiry() time.Duration {
	return time.Duration(envIntOrDefault("EMAIL_VERIFICATION_EXPIRY_MINUTES", defaultVerificationExpiryMinutes)) * time.Minute
}

// VerificationURL builds the link that the user can click to verify their
// email. Set EMAIL_VERIFICATION_LINK to a frontend verification page to send
// users there. When it is empty, the backend verification endpoint is used.
func VerificationURL(token string) string {
	if base := strings.TrimSpace(os.Getenv("EMAIL_VERIFICATION_LINK")); base != "" {
		if parsed, err := url.Parse(base); err == nil {
			query := parsed.Query()
			query.Set("token", token)
			parsed.RawQuery = query.Encode()
			return parsed.String()
		}

		separator := "?"
		if strings.Contains(base, "?") {
			separator = "&"
		}

		return base + separator + "token=" + url.QueryEscape(token)
	}

	return fmt.Sprintf("%s/verify/email/?token=%s", BackendURL(), url.QueryEscape(token))
}

// CreateVerificationChallenge generates a new code and token for the user,
// replacing any previous challenge.
func CreateVerificationChallenge(db *gorm.DB, username string, email string) (*VerificationChallenge, error) {
	code, err := generateNumericCode(verificationCodeDigits)
	if err != nil {
		return nil, err
	}

	token, err := generateVerificationToken()
	if err != nil {
		return nil, err
	}

	expiry := verificationExpiry()

	record := models.EmailVerification{}
	tx := db.Where("username = ?", username).First(&record)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}

	if tx.Error == gorm.ErrRecordNotFound {
		record = models.EmailVerification{Username: username}
	}

	record.Email = email
	record.CodeHash = hashVerificationValue("code:" + code)
	record.TokenHash = hashVerificationValue("token:" + token)
	record.ExpiresAt = time.Now().Add(expiry)
	record.UsedAt = nil
	record.Attempts = 0

	if record.ID == 0 {
		tx = db.Create(&record)
	} else {
		tx = db.Save(&record)
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &VerificationChallenge{
		Record:    &record,
		Code:      code,
		Token:     token,
		Link:      VerificationURL(token),
		ExpiresIn: expiry,
	}, nil
}

// SendAccountVerificationEmail creates a challenge and emails it to the user.
func SendAccountVerificationEmail(db *gorm.DB, username string, name string, email string, role string) error {
	if !EmailVerificationEnabled() {
		return ErrEmailNotConfigured
	}

	if !IsEmailConfigured() {
		return ErrEmailNotConfigured
	}

	challenge, err := CreateVerificationChallenge(db, username, email)
	if err != nil {
		return err
	}

	return SendVerificationEmail(email, name, username, challenge.Code, challenge.Link, role)
}

// ConsumeEmailVerification validates a token, or a username and code pair, and
// marks the challenge as used. It returns the verified username.
func ConsumeEmailVerification(db *gorm.DB, token string, username string, code string) (string, error) {
	token = strings.TrimSpace(token)
	username = strings.TrimSpace(username)
	code = strings.TrimSpace(code)

	if token == "" && (username == "" || code == "") {
		return "", ErrVerificationMissingParams
	}

	record := models.EmailVerification{}
	tx := db

	if token != "" {
		tx = tx.Where("token_hash = ?", hashVerificationValue("token:"+token))
	} else {
		tx = tx.Where("username = ?", username)
	}

	tx = tx.Order("id DESC").First(&record)
	if tx.Error == gorm.ErrRecordNotFound {
		return "", ErrVerificationNotFound
	}
	if tx.Error != nil {
		return "", tx.Error
	}

	if record.UsedAt != nil {
		return "", ErrVerificationAlreadyUsed
	}

	if time.Now().After(record.ExpiresAt) {
		return "", ErrVerificationExpired
	}

	if record.Attempts >= uint(VerificationMaxAttempts()) {
		return "", ErrVerificationAttemptsExceeded
	}

	if token == "" {
		expectedHash := hashVerificationValue("code:" + code)
		if !hmac.Equal([]byte(record.CodeHash), []byte(expectedHash)) {
			attempts := record.Attempts + 1
			if err := db.Model(&models.EmailVerification{}).
				Where("id = ?", record.ID).
				Update("attempts", attempts).Error; err != nil {
				return "", err
			}

			if attempts >= uint(VerificationMaxAttempts()) {
				return "", ErrVerificationAttemptsExceeded
			}

			return "", ErrVerificationCodeInvalid
		}
	}

	now := time.Now()
	tx = db.Model(&models.EmailVerification{}).
		Where("id = ? AND used_at IS NULL", record.ID).
		Updates(map[string]any{
			"used_at":  now,
			"attempts": record.Attempts,
		})

	if tx.Error != nil {
		return "", tx.Error
	}

	if tx.RowsAffected == 0 {
		return "", ErrVerificationAlreadyUsed
	}

	return record.Username, nil
}

// MarkUserEmailVerified marks the student or mentor account as verified.
func MarkUserEmailVerified(db *gorm.DB, username string) (*VerifiedUser, error) {
	student := models.Student{}
	tx := db.Where("username = ?", username).First(&student)
	if tx.Error == nil {
		if !student.EmailVerified {
			if err := db.Model(&models.Student{}).
				Where("id = ?", student.ID).
				Update("email_verified", true).Error; err != nil {
				return nil, err
			}
		}

		return &VerifiedUser{
			Username: student.Username,
			Name:     student.Name,
			Email:    student.Email,
			Role:     "student",
		}, nil
	}
	if tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}

	mentor := models.Mentor{}
	tx = db.Where("username = ?", username).First(&mentor)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, ErrVerificationNotFound
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	if !mentor.EmailVerified {
		if err := db.Model(&models.Mentor{}).
			Where("id = ?", mentor.ID).
			Update("email_verified", true).Error; err != nil {
			return nil, err
		}
	}

	return &VerifiedUser{
		Username: mentor.Username,
		Name:     mentor.Name,
		Email:    mentor.Email,
		Role:     "mentor",
	}, nil
}

func verificationSecret() []byte {
	secret := strings.TrimSpace(os.Getenv("EMAIL_VERIFICATION_SECRET"))
	if secret == "" {
		secret = strings.TrimSpace(os.Getenv("JWT_SECRET_KEY"))
	}
	if secret == "" {
		secret = "kwoc-email-verification-development-secret"
	}

	return []byte(secret)
}

func hashVerificationValue(value string) string {
	mac := hmac.New(sha256.New, verificationSecret())
	_, _ = mac.Write([]byte(value))

	return hex.EncodeToString(mac.Sum(nil))
}

func generateNumericCode(digits int) (string, error) {
	limit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	value, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", digits, value.Int64()), nil
}

func generateVerificationToken() (string, error) {
	buffer := make([]byte, verificationTokenBytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
