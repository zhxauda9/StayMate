package e2e_test

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	baseURL          = "http://localhost:8080"
	chromeDriverPath = "C:\\ChromeDriver\\chromedriver-win64\\chromedriver.exe"
	webDriverURL     = "http://localhost:9515/wd/hub"
)

func TestRegisterLoginAndGetProfileSelenium(t *testing.T) {
	// Start a WebDriver service
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 9515)
	if err != nil {
		t.Fatalf("Failed to start ChromeDriver service: %v", err)
	}
	defer service.Stop()

	// Set up WebDriver capabilities
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{}
	caps.AddChrome(chromeCaps)

	// Connect to the ChromeDriver service
	wd, err := selenium.NewRemote(caps, webDriverURL)
	if err != nil {
		t.Fatalf("Failed to create WebDriver session: %v", err)
	}
	defer wd.Quit()

	// Step 1: Register the user
	t.Log("Registering user")
	if err := wd.Get(baseURL + "/register"); err != nil {
		t.Fatalf("Failed to navigate to register page: %v", err)
	}

	registerEmail := "test@example.com"
	registerPassword := "Aa123456$"
	registerName := "Test User"

	if err := fillInputField(wd, "email", registerEmail); err != nil {
		t.Fatalf("Failed to fill email field: %v", err)
	}
	if err := fillInputField(wd, "password", registerPassword); err != nil {
		t.Fatalf("Failed to fill password field: %v", err)
	}
	if err := fillInputField(wd, "status", registerPassword); err != nil {
		t.Fatalf("Failed to fill password confirmation field: %v", err)
	}
	if err := fillInputField(wd, "name", registerName); err != nil {
		t.Fatalf("Failed to fill name field: %v", err)
	}

	if err := clickButton(wd, "register-button"); err != nil {
		t.Fatalf("Failed to click register button: %v", err)
	}

	time.Sleep(5 * time.Second) // Allow time for registration to complete

	// Step 2: Log in the user
	t.Log("Logging in user")
	if err := wd.Get(baseURL + "/login"); err != nil {
		t.Fatalf("Failed to navigate to login page: %v", err)
	}

	if err := fillInputField(wd, "email", registerEmail); err != nil {
		t.Fatalf("Failed to fill email field: %v", err)
	}
	if err := fillInputField(wd, "password", registerPassword); err != nil {
		t.Fatalf("Failed to fill password field: %v", err)
	}

	if err := clickButton(wd, "login-submit"); err != nil {
		t.Fatalf("Failed to click login button: %v", err)
	}

	time.Sleep(2 * time.Second) // Allow time for login to complete

	// Extract token from cookies
	cookies, err := wd.GetCookies()
	if err != nil {
		t.Fatalf("Failed to retrieve cookies: %v", err)
	}
	token := extractCookieValue(cookies, "Authorization")
	if token == "" {
		t.Fatalf("Authorization token not found")
	}
	t.Logf("Token extracted: %s", token)

	// // Step 3: Get Profile
	// t.Log("Fetching user profile")
	// if err := wd.Get(baseURL + "/profile"); err != nil {
	// 	t.Fatalf("Failed to navigate to profile page: %v", err)
	// }

	if _, err := wd.FindElement(selenium.ByClassName, "profile-container"); err != nil {
		t.Fatalf("Failed to find profile content: %v", err)
	}

	t.Log("Test completed successfully")
}

func fillInputField(wd selenium.WebDriver, fieldID, value string) error {
	field, err := wd.FindElement(selenium.ByID, fieldID)
	if err != nil {
		return err
	}
	return field.SendKeys(value)
}

func clickButton(wd selenium.WebDriver, buttonID string) error {
	button, err := wd.FindElement(selenium.ByID, buttonID)
	if err != nil {
		return err
	}
	return button.Click()
}

func extractCookieValue(cookies []selenium.Cookie, name string) string {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}
