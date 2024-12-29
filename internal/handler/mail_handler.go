package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"text/template"

	"github.com/zhxauda9/StayMate/internal/config"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
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
		l.Log.Error().Err(err).Msg("Error loading template")
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Executing template
	tmpl.Execute(w, data)
}

func (h *mailHandler) SendMailHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Limit to 10 MB
		l.Log.Error().Err(err).Msg("Unable to parse form")
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ensure r.MultipartForm is not nil
	if r.MultipartForm == nil {
		l.Log.Error().Msg("Invalid form data")
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get and validate email
	emails, emailExists := getField(r, "emails")
	if !emailExists {
		l.Log.Error().Msg("Email field is missing")
		http.Error(w, "Email field is missing", http.StatusBadRequest)
		return
	}

	// Get and validate subject
	subject, subjectExists := getField(r, "subject")
	if !subjectExists {
		l.Log.Error().Msg("Subject field is missing")
		http.Error(w, "Subject field is missing", http.StatusBadRequest)
		return
	}

	// Get and validate message
	message, messageExists := getField(r, "message")
	if !messageExists {
		l.Log.Error().Msg("Message field is missing")
		http.Error(w, "Message field is missing", http.StatusBadRequest)
		return
	}

	err := h.mailService.Send(emails, subject[0], message[0], "", "", nil)
	if err != nil {
		l.Log.Error().Err(err).Msg("Failed sending email")
		http.Error(w, "Failed sending email", http.StatusInternalServerError)
	}

	l.Log.Info().Strs("emails", emails).Str("subject", subject[0]).Str("message", message[0]).Msg("Mails sended successfully")
	// Respond to the client
	w.Write([]byte("Email sended successfully"))
}

// testing don't work yet
func (h *mailHandler) SendMailFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	if err := r.ParseMultipartForm(20 << 20); err != nil { // Limit to 20 MB
		l.Log.Error().Err(err).Msg("Unable to parse form")
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// File parsing
	file, header, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "File not provided", http.StatusNotFound)
			return
		} else {
			l.Log.Error().Err(err).Msg("Unable to get file from form")
			http.Error(w, "Unable to retrieve the file or the file not provided", http.StatusBadRequest)
			return
		}
	}
	defer file.Close()

	mimeType := r.MultipartForm.File["file"][0].Header.Get("Content-Type")

	// Validating mimeType of file
	if !config.AvailiableMimeTypes[mimeType] {
		l.Log.Error().Str("MimeType", mimeType).Msg("Invalid MimeType type")
		http.Error(w, fmt.Sprintf("Invalid MimeType: %v", mimeType), http.StatusBadRequest)
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		l.Log.Error().Err(err).Msg("Unable to read file")
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Getting list of emails
	emails, emailExists := getField(r, "emails")
	if !emailExists {
		l.Log.Error().Msg("Emails field is missing")
		http.Error(w, "emails field is missing", http.StatusBadRequest)
		return
	}

	// Get and validate subject
	subject, subjectExists := getField(r, "subject")
	if !subjectExists {
		l.Log.Error().Msg("Subject field is missing")
		http.Error(w, "Subject field is missing", http.StatusBadRequest)
		return
	}

	// Get and validate message
	message, messageExists := getField(r, "message")
	if !messageExists {
		l.Log.Error().Msg("Message field is missing")
		http.Error(w, "Message field is missing", http.StatusBadRequest)
		return
	}

	// Validate email addresses
	l.Log.Debug().Int("number", len(emails)).Strs("emails", emails).Msg("Validating emails from request")
	for _, email := range emails {
		if _, err := mail.ParseAddress(email); err != nil {
			l.Log.Error().Str("email", email).Err(err).Msg("Invalid email address")
			http.Error(w, fmt.Sprintf("Invalid email address: %s", email), http.StatusBadRequest)
			return
		}
	}

	l.Log.Debug().Int("number", len(emails)).Strs("emails", emails).Msg("Mails validated successfully")

	// Send(mails []string, subject, message, filename, mimeType string, filedata []byte) error
	l.Log.Debug().Str("filename", header.Filename).Str("MimeType", mimeType).Int("Data length", len(fileData)).Msg("File info")
	err = h.mailService.Send(emails, subject[0], message[0], header.Filename, mimeType, fileData)
	if err != nil {
		l.Log.Error().Err(err).Msg("Failed to send email")
		http.Error(w, "Failed to send email", http.StatusBadRequest)
		return
	}

	l.Log.Info().Strs("emails", emails).Str("subject", subject[0]).Str("message", message[0]).Msg("Mails sended successfully")
	w.Write([]byte("Email sended successfully"))
}

func getField(r *http.Request, field string) ([]string, bool) {
	values, exists := r.MultipartForm.Value[field]
	if !exists || len(values) == 0 {
		return nil, false
	}
	return values, true
}
