package validation_email

import (
	"encoding/json"
	"fmt"
	"net/mail"
)

type EmailRequest struct {
	Email string `json:"email"`
}

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func ValidateEmail(email string) (bool, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return false, &ValidationError{
			Code:    "invalid_format",
			Message: "Invalid email format",
		}
	}

	// Try primary API first
	body, err := CallPrimaryAPI(email)
	if err != nil {
		// Fallback to secondary API
		body, err = CallSecondaryAPI(email)
		if err != nil {
			return false, err
		}
		return parseSecondaryAPIResponse(body)
	}
	return parsePrimaryAPIResponse(body)
}

// Helper for primary API response
func parsePrimaryAPIResponse(body []byte) (bool, error) {
	var resp struct {
		Result struct {
			ValidationDetails struct {
				FormatValid bool `json:"format_valid"`
				Disposable  bool `json:"disposable"`
				SMTPCheck   bool `json:"smtp_check"`
			} `json:"validation_details"`
			Status string  `json:"status"`
			Score  float64 `json:"score"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, fmt.Errorf("invalid primary API response")
	}
	// Validate: must be valid format, not disposable, and exist
	if !resp.Result.ValidationDetails.FormatValid ||
		resp.Result.ValidationDetails.Disposable ||
		!resp.Result.ValidationDetails.SMTPCheck {
		return false, nil
	}
	return true, nil
}

// Helper for secondary API response
func parseSecondaryAPIResponse(body []byte) (bool, error) {
	var resp struct {
		Result struct {
			IsValid      bool `json:"is_valid"`
			IsDisposable bool `json:"is_disposable"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, fmt.Errorf("invalid secondary API response")
	}
	if !resp.Result.IsValid || resp.Result.IsDisposable {
		return false, nil
	}
	return true, nil
}
