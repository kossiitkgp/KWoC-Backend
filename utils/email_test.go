package utils

import (
	"bufio"
	"encoding/base64"
	"errors"
	"net"
	"net/mail"
	"strings"
	"testing"
	"time"
)

func TestGenerateNumericCode(t *testing.T) {
	code, err := generateNumericCode(verificationCodeDigits)
	if err != nil {
		t.Fatalf("Expected no error. Got %v.", err)
	}

	if len(code) != verificationCodeDigits {
		t.Errorf("Expected a %d digit code. Got `%s`.", verificationCodeDigits, code)
	}

	for _, char := range code {
		if char < '0' || char > '9' {
			t.Errorf("Expected only digits. Got `%s`.", code)
			break
		}
	}
}

func TestHashVerificationValue(t *testing.T) {
	t.Setenv("EMAIL_VERIFICATION_SECRET", "test-secret")

	first := hashVerificationValue("code:123456")
	second := hashVerificationValue("code:123456")
	third := hashVerificationValue("token:123456")

	if first != second {
		t.Error("Expected the same value to produce the same hash.")
	}

	if first == third {
		t.Error("Expected different values to produce different hashes.")
	}

	if len(first) != 64 {
		t.Errorf("Expected a 64 character SHA-256 hex hash. Got `%s`.", first)
	}
}

func TestVerificationURL(t *testing.T) {
	t.Setenv("EMAIL_VERIFICATION_LINK", "https://kwoc.example.com/verify-email?source=email")
	t.Setenv("BACKEND_PUBLIC_URL", "https://api.kwoc.example.com")

	link := VerificationURL("abc 123")
	if link != "https://kwoc.example.com/verify-email?source=email&token=abc+123" {
		t.Errorf("Unexpected verification link `%s`.", link)
	}

	t.Setenv("EMAIL_VERIFICATION_LINK", "")
	link = VerificationURL("abc123")
	if link != "https://api.kwoc.example.com/verify/email/?token=abc123" {
		t.Errorf("Unexpected verification link `%s`.", link)
	}
}

func TestBuildMailMessage(t *testing.T) {
	from := &mail.Address{Address: "no-reply@example.com"}
	to := &mail.Address{Address: "student@example.com"}

	message, err := buildMailMessage(from, "KWoC", to, Mail{
		To:      to.Address,
		Subject: "Verify your email\r\nBcc: attacker@example.com",
		Text:    "Your code is 123456.",
		HTML:    "<p>Your code is <strong>123456</strong>.</p>",
	})
	if err != nil {
		t.Fatalf("Expected no error. Got %v.", err)
	}

	raw := string(message)

	if strings.Contains(raw, "\r\nBcc:") {
		t.Error("Expected header injection to be removed.")
	}

	expectedHeaders := []string{
		"From: \"KWoC\" <no-reply@example.com>",
		"To: <student@example.com>",
		"MIME-Version: 1.0",
		"multipart/alternative",
	}
	for _, header := range expectedHeaders {
		if !strings.Contains(raw, header) {
			t.Errorf("Expected message to contain `%s`.", header)
		}
	}

	parts := strings.SplitN(raw, "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatal("Expected message headers and body.")
	}

	if !strings.Contains(parts[1], "text/plain") || !strings.Contains(parts[1], "text/html") {
		t.Error("Expected both text and HTML parts.")
	}
}

func TestBuildMailMessageTextOnly(t *testing.T) {
	from := &mail.Address{Address: "no-reply@example.com"}
	to := &mail.Address{Address: "student@example.com"}

	message, err := buildMailMessage(from, "KWoC", to, Mail{
		To:      to.Address,
		Subject: "Hello",
		Text:    "Plain text body.",
	})
	if err != nil {
		t.Fatalf("Expected no error. Got %v.", err)
	}

	parts := strings.SplitN(string(message), "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatal("Expected message headers and body.")
	}

	body := strings.ReplaceAll(parts[1], "\r\n", "")
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		t.Fatalf("Expected a base64 body. Got %v.", err)
	}

	if string(decoded) != "Plain text body." {
		t.Errorf("Unexpected body `%s`.", string(decoded))
	}
}

func TestMailerNotConfigured(t *testing.T) {
	t.Setenv("SMTP_HOST", "")
	t.Setenv("EMAIL_FROM_ADDRESS", "")
	t.Setenv("SMTP_USERNAME", "")

	mailer := NewMailerFromEnv()
	if mailer.Configured() {
		t.Error("Expected mailer to be unconfigured.")
	}

	err := mailer.Send(Mail{To: "student@example.com", Subject: "Hello", Text: "Hi"})
	if !errors.Is(err, ErrEmailNotConfigured) {
		t.Errorf("Expected ErrEmailNotConfigured. Got %v.", err)
	}
}

func TestMailerInvalidRecipient(t *testing.T) {
	mailer := &Mailer{
		Host:      "smtp.example.com",
		FromEmail: "no-reply@example.com",
	}

	err := mailer.Send(Mail{To: "not an email", Subject: "Hello", Text: "Hi"})
	if !errors.Is(err, ErrInvalidEmailAddress) {
		t.Errorf("Expected ErrInvalidEmailAddress. Got %v.", err)
	}
}

func TestMailerSend(t *testing.T) {
	host, port, messages, stop := startFakeSMTPServer(t)
	defer stop()

	t.Setenv("SMTP_HOST", host)
	t.Setenv("SMTP_PORT", port)
	t.Setenv("SMTP_ENCRYPTION", "none")
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_PASSWORD", "")
	t.Setenv("EMAIL_FROM_ADDRESS", "no-reply@example.com")
	t.Setenv("EMAIL_FROM_NAME", "KWoC")

	err := SendEmail(Mail{
		To:      "student@example.com",
		Subject: "Verify your email",
		Text:    "Your code is 123456.",
		HTML:    "<p>Your code is <strong>123456</strong>.</p>",
	})
	if err != nil {
		t.Fatalf("Expected no error. Got %v.", err)
	}

	select {
	case message := <-messages:
		if !strings.Contains(message, "Subject: Verify your email") {
			t.Error("Expected the subject in the captured email.")
		}

		if !strings.Contains(message, "text/plain") || !strings.Contains(message, "text/html") {
			t.Error("Expected both text and HTML parts in the captured email.")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for the email.")
	}
}

// startFakeSMTPServer starts a minimal SMTP server for tests. It returns the
// host and port to connect to, a channel that receives the raw email data, and
// a cleanup function.
func startFakeSMTPServer(t *testing.T) (string, string, chan string, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error starting fake SMTP server: %v", err)
	}

	messages := make(chan string, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer func() {
			_ = conn.Close()
		}()

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		writeLine := func(line string) bool {
			if _, err := writer.WriteString(line + "\r\n"); err != nil {
				return false
			}

			return writer.Flush() == nil
		}

		if !writeLine("220 localhost ESMTP") {
			return
		}

		var data strings.Builder
		inData := false

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimRight(line, "\r\n")
			upper := strings.ToUpper(line)

			if inData {
				if line == "." {
					inData = false
					if !writeLine("250 OK") {
						return
					}
					continue
				}

				data.WriteString(line)
				data.WriteString("\r\n")
				continue
			}

			switch {
			case strings.HasPrefix(upper, "EHLO"), strings.HasPrefix(upper, "HELO"):
				if !writeLine("250-localhost\r\n250 OK") {
					return
				}
			case strings.HasPrefix(upper, "MAIL FROM"),
				strings.HasPrefix(upper, "RCPT TO"),
				upper == "RSET":
				if !writeLine("250 OK") {
					return
				}
			case upper == "DATA":
				if !writeLine("354 End data with <CR><LF>.<CR><LF>") {
					return
				}
				inData = true
			case upper == "QUIT":
				_ = writeLine("221 Bye")
				messages <- data.String()
				return
			default:
				if !writeLine("250 OK") {
					return
				}
			}
		}
	}()

	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("Error getting fake SMTP server address: %v", err)
	}

	stop := func() {
		_ = listener.Close()
	}

	return host, port, messages, stop
}
