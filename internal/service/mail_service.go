package service

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/go-mail/mail/v2"
)

var (
	ErrNoEmails   = errors.New("no emails provided")
	ErrNoMimeType = errors.New("mimeType of file not provided")
	ErrNoFileData = errors.New("data of file not provided")
)

type MailServiceImpl interface {
	Send(mails []string, subject, message, filename, mimeType string, filedata []byte) error
}

type mailService struct {
	smtpHost string
	smtpPort int
	email    string
	password string
	smptAddr string
	dialer   *mail.Dialer
}

func NewMailService(smtpHost, smtpPort, email, password string) (*mailService, error) {
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, errors.New("Invalid port. Error: " + err.Error())
	}

	d := mail.NewDialer(smtpHost, port, email, password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return &mailService{
		smtpHost: smtpHost,
		smtpPort: port,
		email:    email,
		password: password,
		smptAddr: smtpHost + ":" + smtpPort,
		dialer:   d,
	}, nil
}

/*
Sends a file to all given emails. Subject and message are optional.
MimeType and filedata are necessary, if filename provided.
Method won't send file if filename not provided
*/
func (s *mailService) Send(mails []string, subject, message, filename, mimeType string, filedata []byte) error {
	if len(mails) == 0 {
		return ErrNoEmails
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.email)
	m.SetHeader("To", mails...)

	if subject != "" {
		m.SetHeader("Subject", subject)
	}
	if message != "" {
		m.SetBody("text/plain", message)
	}

	// Attaching file if filename provided
	if filename != "" {
		if mimeType == "" {
			return ErrNoMimeType
		}
		if len(filedata) == 0 {
			return ErrNoFileData
		}

		m.Attach(filename,
			mail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(filedata) // Write the file data directly to the writer
				return err
			}),
			mail.SetHeader(map[string][]string{
				"Content-Type":        {fmt.Sprintf("%s; name=\"%s\"", mimeType, filename)},
				"Content-Disposition": {fmt.Sprintf("attachment; filename=\"%s\"", filename)},
			}),
		)
	}

	// Sending the email
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
