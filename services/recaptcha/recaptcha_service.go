package recaptcha

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "time"
)

type VerifyResponse struct {
    Success     bool     `json:"success"`
    Score       float64  `json:"score"`
    Action      string   `json:"action"`
    ChallengeTS string   `json:"challenge_ts"`
    Hostname    string   `json:"hostname"`
    ErrorCodes  []string `json:"error-codes"`
}

var siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

// Verify verifies a reCAPTCHA v3 token using the secret from env.
// expectedAction: action name you executed on the client (e.g. "contact_submit").
// minScore: recommended threshold (e.g. 0.5). Set to 0 to skip score check.
func Verify(token string, expectedAction string, minScore float64, remoteIP string) (bool, float64, error) {
    secret := os.Getenv("RECAPTCHA_SECRET_KEY")
    if secret == "" {
        return false, 0, fmt.Errorf("recaptcha secret not configured")
    }
    if token == "" {
        return false, 0, fmt.Errorf("recaptcha token empty")
    }

    form := url.Values{}
    form.Add("secret", secret)
    form.Add("response", token)
    if remoteIP != "" {
        form.Add("remoteip", remoteIP)
    }

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.PostForm(siteVerifyURL, form)
    if err != nil {
        return false, 0, err
    }
    defer resp.Body.Close()

    var vr VerifyResponse
    if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
        return false, 0, fmt.Errorf("invalid recaptcha response: %w", err)
    }

    // Basic checks: success, optional action match, optional score threshold
    if !vr.Success {
        return false, vr.Score, fmt.Errorf("recaptcha verification failed: %v", vr.ErrorCodes)
    }
    if expectedAction != "" && vr.Action != expectedAction {
        return false, vr.Score, fmt.Errorf("recaptcha action mismatch (got=%s expected=%s)", vr.Action, expectedAction)
    }
    if minScore > 0 && vr.Score < minScore {
        return false, vr.Score, fmt.Errorf("recaptcha score too low: %f < %f", vr.Score, minScore)
    }
    return true, vr.Score, nil
}