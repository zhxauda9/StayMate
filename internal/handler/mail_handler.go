package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"text/template"

	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type mailHandler struct {
	mailService service.MailServiceImpl
	userService service.UserService
}

func NewMailHandler(mailService service.MailServiceImpl, userService service.UserService) *mailHandler {
	return &mailHandler{mailService: mailService, userService: userService}
}

// Serves template
func (h *mailHandler) ServeMail(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.URL.Query().Get("userId")

	// if parametr doesn't provided
	if userIdParam == "" {
		userIdParam = "1"
	}

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	//Getting user
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Could not get users", http.StatusInternalServerError)
		return
	}

	data := struct {
		SelectedUserID int
		Users          []models.User
	}{
		SelectedUserID: userId,
		Users:          users,
	}

	// Preparing template
	tmpl, err := template.ParseFiles("web/templates/send-email.html")
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Error loading template")
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Executing template
	tmpl.Execute(w, data)
}

// testing don't work yet
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
