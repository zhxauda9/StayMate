package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"

	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
)

type mailHandler struct {
	mailService service.MailServiceImpl
}

func NewMailHandler(mailService service.MailServiceImpl) *mailHandler {
	return &mailHandler{mailService: mailService}
}

func (h *mailHandler) SendFileMailHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Unable to get file from form")
		http.Error(w, "Unable to retrieve the file or the file not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := r.MultipartForm.File["file"][0].Header.Get("Content-Type")

	fileData, err := io.ReadAll(file)
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Unable to read file")
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Getting list of mails to send the file
	emails := r.MultipartForm.Value["emails"]

	// Validate email addresses
	myLogger.Log.Debug().Int("number", len(emails)).Strs("emails", emails).Msg("Validating emails from request")
	var emailList []string
	for _, email := range emails {
		if _, err := mail.ParseAddress(email); err != nil {
			myLogger.Log.Error().Str("email", email).Err(err).Msg("Invalid email address")
			http.Error(w, fmt.Sprintf("Invalid email address: %s", email), http.StatusBadRequest)
			return
		}
		emailList = append(emailList, email)
	}
	myLogger.Log.Debug().Int("number", len(emailList)).Strs("emails", emailList).Msg("Mails validated successfully")

	// Send(mails []string, subject, message, filename, mimeType string, filedata []byte) error
	err = h.mailService.Send(emailList, "Staymate subject", "Hello from StayMate", header.Filename, mimeType, fileData)
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Could not send email")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}
