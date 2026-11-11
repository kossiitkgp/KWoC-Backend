package models

import (
	"time"

	"gorm.io/gorm"
)

// `email_verifications` table model

// EmailVerification stores the active email verification challenge for a user.
// Only the hashes of the code and the token are stored.
type EmailVerification struct {
	gorm.Model

	// Username of the student or mentor being verified. A user can only have
	// one active challenge at a time.
	Username string `gorm:"size:255;uniqueIndex;column:username"`

	// Email address the challenge was sent to.
	Email string `gorm:"size:255;column:email"`

	// SHA-256 HMAC of the verification code.
	CodeHash string `gorm:"size:64;column:code_hash"`

	// SHA-256 HMAC of the verification token (used in the verification link).
	TokenHash string `gorm:"size:64;uniqueIndex;column:token_hash"`

	ExpiresAt time.Time `gorm:"column:expires_at"`

	// UsedAt is set once the challenge has been successfully consumed.
	UsedAt *time.Time `gorm:"column:used_at"`

	// Attempts is the number of failed code submissions for this challenge.
	Attempts uint `gorm:"column:attempts"`
}
