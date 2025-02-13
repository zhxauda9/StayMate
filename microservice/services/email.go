package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	_ "strings"
)

type EmailSettings struct {
	Host string
	Port string
	User string
	Pass string
}

// SendReceiptEmail sends an email with the PDF receipt as an attachment.
func SendReceiptEmail(settings EmailSettings, to, subject string, body string, attachmentName string, attachmentBytes []byte) error {
	auth := smtp.PlainAuth("", settings.User, settings.Pass, settings.Host)

	// MIME headers for the email
	header := make(map[string]string)
	header["From"] = settings.User
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `multipart/mixed; boundary="boundary123"`

	// Create email body with attachment
	var msg bytes.Buffer
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n--boundary123\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	msg.WriteString(body + "\r\n")
	msg.WriteString("\r\n--boundary123\r\n")
	msg.WriteString("Content-Type: application/pdf\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attachmentName))

	// Base64 encode the attachment and split it into lines of 76 characters
	encoded := base64.StdEncoding.EncodeToString(attachmentBytes)
	lines := splitByLength(encoded, 76)
	for _, line := range lines {
		msg.WriteString(line + "\r\n")
	}

	msg.WriteString("--boundary123--")

	// Send the email
	address := fmt.Sprintf("%s:%s", settings.Host, settings.Port)
	return smtp.SendMail(address, auth, settings.User, []string{to}, msg.Bytes())
}

// Helper function to split a string into chunks of the specified length
func splitByLength(s string, length int) []string {
	var result []string
	for len(s) > 0 {
		if len(s) > length {
			result = append(result, s[:length])
			s = s[length:]
		} else {
			result = append(result, s)
			s = ""
		}
	}
	return result
}
