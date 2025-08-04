package validation_email

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func CallPrimaryAPI(email string) ([]byte, error) {
	apiURL := os.Getenv("PRIMARY_EMAIL_VERIFY_URL")
	apiKey := os.Getenv("PRIMARY_EMAIL_VERIFY_KEY")
	url := fmt.Sprintf("%s?api_key=%s&email_address=%s", apiURL, apiKey, email)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, fmt.Errorf("primary API failed")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
