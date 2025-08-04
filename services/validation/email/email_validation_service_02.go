package validation_email

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func CallSecondaryAPI(email string) ([]byte, error) {
	apiURL := os.Getenv("SECONDARY_EMAIL_VERIFY_URL")
	apiKey := os.Getenv("SECONDARY_EMAIL_VERIFY_KEY")
	url := fmt.Sprintf("%s?email=%s", apiURL, email)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, fmt.Errorf("secondary API failed")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
