package utils

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultSMTPPort       = "587"
	defaultMailTimeout    = 15 * time.Second
	defaultFromName       = "KWoC"
	defaultSupportEmail   = "kwoc@kossiitkgp.org"
	emailEncryptionNone   = "none"
	emailEncryptionTLS    = "tls"
	emailEncryptionSSL    = "ssl"
	emailEncryptionStrtls = "starttls"
)

var (
	ErrEmailNotConfigured   = errors.New("email is not configured")
	ErrEmailRecipientEmpty  = errors.New("email recipient is empty")
	ErrInvalidEmailAddress  = errors.New("invalid email address")
	ErrEmailStartTlsMissing = errors.New("SMTP server does not support STARTTLS")
)

// Mail is a single outgoing email.
type Mail struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

// Mailer sends email through an SMTP server.
type Mailer struct {
	Host       string
	Port       string
	Username   string
	Password   string
	FromEmail  string
	FromName   string
	Encryption string
	Timeout    time.Duration
}

// NewMailerFromEnv builds a Mailer from the SMTP environment variables.
func NewMailerFromEnv() *Mailer {
	fromEmail := strings.TrimSpace(os.Getenv("EMAIL_FROM_ADDRESS"))
	if fromEmail == "" {
		fromEmail = strings.TrimSpace(os.Getenv("SMTP_USERNAME"))
	}

	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if port == "" {
		port = defaultSMTPPort
	}

	encryption := strings.ToLower(strings.TrimSpace(os.Getenv("SMTP_ENCRYPTION")))
	if encryption == "" {
		encryption = emailEncryptionStrtls
	}

	return &Mailer{
		Host:       strings.TrimSpace(os.Getenv("SMTP_HOST")),
		Port:       port,
		Username:   os.Getenv("SMTP_USERNAME"),
		Password:   os.Getenv("SMTP_PASSWORD"),
		FromEmail:  fromEmail,
		FromName:   envOrDefault("EMAIL_FROM_NAME", defaultFromName),
		Encryption: encryption,
		Timeout:    defaultMailTimeout,
	}
}

// IsEmailConfigured reports whether enough SMTP configuration is present to
// attempt sending an email.
func IsEmailConfigured() bool {
	return NewMailerFromEnv().Configured()
}

// Configured reports whether the mailer has the minimum required settings.
func (m *Mailer) Configured() bool {
	return m != nil && m.Host != "" && m.FromEmail != ""
}

// Send delivers a single email. The call is synchronous and bounded by the
// mailer timeout.
func (m *Mailer) Send(email Mail) error {
	if !m.Configured() {
		return ErrEmailNotConfigured
	}

	if strings.TrimSpace(email.To) == "" {
		return ErrEmailRecipientEmpty
	}

	to, err := parseEmailAddress(email.To)
	if err != nil {
		return ErrInvalidEmailAddress
	}

	from, err := parseEmailAddress(m.FromEmail)
	if err != nil {
		return ErrInvalidEmailAddress
	}

	message, err := buildMailMessage(from, m.FromName, to, email)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(m.Host, m.Port)
	dialer := &net.Dialer{Timeout: m.Timeout}

	var conn net.Conn
	if m.Encryption == emailEncryptionSSL || m.Encryption == emailEncryptionTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
			ServerName: m.Host,
			MinVersion: tls.VersionTLS12,
		})
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return err
	}

	// Bound the whole SMTP conversation, not just the TCP handshake.
	if err := conn.SetDeadline(time.Now().Add(m.Timeout)); err != nil {
		_ = conn.Close()
		return err
	}

	client, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	if m.Encryption == emailEncryptionStrtls {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{
				ServerName: m.Host,
				MinVersion: tls.VersionTLS12,
			}); err != nil {
				return err
			}
		} else if m.Username != "" {
			// Do not send credentials over a plaintext connection.
			return ErrEmailStartTlsMissing
		}
	}

	if m.Username != "" {
		auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(from.Address); err != nil {
		return err
	}

	if err := client.Rcpt(to.Address); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

// SendEmail is a convenience wrapper around Mailer.Send.
func SendEmail(email Mail) error {
	return NewMailerFromEnv().Send(email)
}

// SupportEmail returns the address users can contact for help.
func SupportEmail() string {
	return envOrDefault("SUPPORT_EMAIL", defaultSupportEmail)
}

// FrontendURL returns the configured frontend URL without a trailing slash.
func FrontendURL() string {
	return strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_URL")), "/")
}

// BackendURL returns the configured public backend URL without a trailing slash.
func BackendURL() string {
	base := strings.TrimRight(strings.TrimSpace(os.Getenv("BACKEND_PUBLIC_URL")), "/")
	if base != "" {
		return base
	}

	return "http://localhost:" + envOrDefault("BACKEND_PORT", "8080")
}

func parseEmailAddress(raw string) (*mail.Address, error) {
	raw = sanitizeEmailHeader(raw)

	address, err := mail.ParseAddress(raw)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func buildMailMessage(from *mail.Address, fromName string, to *mail.Address, email Mail) ([]byte, error) {
	if fromName != "" {
		from.Name = fromName
	}

	subject := sanitizeEmailHeader(email.Subject)
	if subject == "" {
		subject = "KWoC"
	}

	var builder strings.Builder
	writeHeader(&builder, "From", from.String())
	writeHeader(&builder, "To", to.String())
	writeHeader(&builder, "Subject", mime.QEncoding.Encode("UTF-8", subject))
	writeHeader(&builder, "Date", time.Now().Format(time.RFC1123Z))
	writeHeader(&builder, "MIME-Version", "1.0")
	writeHeader(&builder, "Auto-Submitted", "auto-generated")

	textBody := email.Text
	htmlBody := email.HTML

	switch {
	case textBody != "" && htmlBody != "":
		boundary := fmt.Sprintf("kwoc-%d", time.Now().UnixNano())
		writeHeader(&builder, "Content-Type", fmt.Sprintf(`multipart/alternative; boundary="%s"`, boundary))
		builder.WriteString("\r\n")

		writeMIMEPart(&builder, boundary, "text/plain; charset=\"UTF-8\"", textBody)
		writeMIMEPart(&builder, boundary, "text/html; charset=\"UTF-8\"", htmlBody)

		fmt.Fprintf(&builder, "--%s--\r\n", boundary)
	case htmlBody != "":
		writeHeader(&builder, "Content-Type", `text/html; charset="UTF-8"`)
		writeHeader(&builder, "Content-Transfer-Encoding", "base64")
		builder.WriteString("\r\n")
		builder.WriteString(encodeBase64Lines(htmlBody))
	default:
		writeHeader(&builder, "Content-Type", `text/plain; charset="UTF-8"`)
		writeHeader(&builder, "Content-Transfer-Encoding", "base64")
		builder.WriteString("\r\n")
		builder.WriteString(encodeBase64Lines(textBody))
	}

	return []byte(builder.String()), nil
}

func writeHeader(builder *strings.Builder, key string, value string) {
	fmt.Fprintf(builder, "%s: %s\r\n", key, sanitizeEmailHeader(value))
}

func writeMIMEPart(builder *strings.Builder, boundary string, contentType string, body string) {
	fmt.Fprintf(builder, "--%s\r\n", boundary)
	writeHeader(builder, "Content-Type", contentType)
	writeHeader(builder, "Content-Transfer-Encoding", "base64")
	builder.WriteString("\r\n")
	builder.WriteString(encodeBase64Lines(body))
	builder.WriteString("\r\n")
}

func encodeBase64Lines(value string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))

	var builder strings.Builder
	for len(encoded) > 76 {
		builder.WriteString(encoded[:76])
		builder.WriteString("\r\n")
		encoded = encoded[76:]
	}
	builder.WriteString(encoded)
	builder.WriteString("\r\n")

	return builder.String()
}

// sanitizeEmailHeader removes characters that could be used for header
// injection.
func sanitizeEmailHeader(value string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(strings.TrimSpace(value))
}

func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func envIntOrDefault(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		log.Warn().Msgf("Invalid value for %s: %s. Using %d.", key, value, fallback)
		return fallback
	}

	return parsed
}
