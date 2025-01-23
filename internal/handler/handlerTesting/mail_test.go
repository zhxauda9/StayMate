package handlertesting

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zhxauda9/StayMate/internal/handler"
	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/models"
)

// Mock MailService
type mockMailService struct {
	mock.Mock
}

func (m *mockMailService) Send(mails []string, subject, message, filename, mimeType string, filedata []byte) error {
	args := m.Called(mails, subject, message, filename, mimeType, filedata)
	return args.Error(0)
}

// Mock UserService
type fakeService struct{}

func (m *fakeService) CreateUser(user models.User) error {
	return nil
}

func (m *fakeService) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

func (m *fakeService) GetUserByEmail(email string) (models.User, error) {
	return models.User{}, nil
}

func (m *fakeService) GetAllUsers(sort string, page int) ([]models.User, error) {
	return nil, nil
}

func (m *fakeService) UpdateUser(id int, user models.User) error {
	return nil
}

func (m *fakeService) DeleteUser(id int) error {
	return nil
}

// Test SendMailHandler Method with valid input
func TestSendMailHandler(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLogger()
	mockMailService := new(mockMailService)
	mockUserService := new(fakeService)
	handler := handler.NewMailHandler(mockMailService, mockUserService)

	// Prepare form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("emails", "test@example.com")
	writer.WriteField("subject", "Test Subject")
	writer.WriteField("message", "Test Message")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/mail/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	// Mocking mail service response
	mockMailService.On("Send", []string{"test@example.com"}, "Test Subject", "Test Message", "", "", []byte(nil)).Return(nil)

	handler.SendMailHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Email sended successfully", rec.Body.String())
	mockMailService.AssertExpectations(t)
}

// Test SendMailHandler with missing email field
func TestSendMailHandler_MissingEmail(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLogger()
	mockMailService := new(mockMailService)
	mockUserService := new(fakeService)
	handler := handler.NewMailHandler(mockMailService, mockUserService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("subject", "Test Subject")
	writer.WriteField("message", "Test Message")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/mail/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.SendMailHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email field is missing")
}

// Test SendMailFileHandler with invalid mime type
func TestSendMailFileHandler_InvalidMimeType(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLogger()
	mockMailService := new(mockMailService)
	mockUserService := new(fakeService)
	handler := handler.NewMailHandler(mockMailService, mockUserService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("emails", "test@example.com")
	writer.WriteField("subject", "Test Subject")
	writer.WriteField("message", "Test Message")
	// Add a file with invalid mime type
	file, _ := writer.CreateFormFile("file", "test.txt")
	file.Write([]byte("Test File Content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/mail/sendfile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.SendMailFileHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid MimeType")
}

// Test SendMailFileHandler with missing file
func TestSendMailFileHandler_MissingFile(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLogger()
	mockMailService := new(mockMailService)
	mockUserService := new(fakeService)
	handler := handler.NewMailHandler(mockMailService, mockUserService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("emails", "test@example.com")
	writer.WriteField("subject", "Test Subject")
	writer.WriteField("message", "Test Message")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/mail/sendfile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.SendMailFileHandler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "File not provided")
}

// Still working with this tests
// Test SendMailFileHandler Method with a real PDF file
// func TestSendMailFileHandler(t *testing.T) {
// 	myLogger.Log = myLogger.NewZeroLogger()
// 	mockMailService := new(mockMailService) // Replace with the actual mock service
// 	mockUserService := new(fakeService)
// 	handler := handler.NewMailHandler(mockMailService, mockUserService)

// 	// Open a real PDF file from the given directory
// 	filePath := "./TESTFILE.pdf" // Specify the path to the test.pdf file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		t.Fatalf("Error opening file: %v", err)
// 	}
// 	defer file.Close()

// 	// Detect the MIME type based on the file content
// 	mimeType := mime.TypeByExtension(".pdf") // Use the file extension to detect MIME type
// 	if mimeType == "" {
// 		t.Fatalf("Failed to detect MIME type")
// 	}

// 	// Prepare multipart form with the real file
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	writer.WriteField("emails", "test@example.com")
// 	writer.WriteField("subject", "Test Subject")
// 	writer.WriteField("message", "Test Message")

// 	// Add the real file (test.pdf) to the multipart form
// 	part, err := writer.CreateFormFile("file", "test.pdf")
// 	if err != nil {
// 		t.Fatalf("Error creating form file: %v", err)
// 	}

// 	// Copy the real file content into the multipart form
// 	_, err = io.Copy(part, file)
// 	if err != nil {
// 		t.Fatalf("Error copying file content: %v", err)
// 	}

// 	writer.Close()

// 	req := httptest.NewRequest(http.MethodPost, "/mail/sendfile", body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	rec := httptest.NewRecorder()
// 	t.Logf("Sending MIME TYPE %v", writer.FormDataContentType())
// 	// Mocking mail service response with the correct MIME type (application/pdf)
// 	mockMailService.On("Send", []string{"test@example.com"}, "Test Subject", "Test Message", "test.pdf", mimeType, mock.Anything).Return(nil)

// 	// Call handler method to send the email
// 	handler.SendMailFileHandler(rec, req)

// 	// Asserting that the email was sent successfully
// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	assert.Equal(t, "Email sended successfully", rec.Body.String())
// 	mockMailService.AssertExpectations(t)
// }
