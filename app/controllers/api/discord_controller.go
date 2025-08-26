package api_controllers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
)

type DiscordController struct{}

func NewDiscordController() *DiscordController {
	return &DiscordController{}
}

func (dc *DiscordController) SendWebhook(w http.ResponseWriter, r *http.Request) {
    // Optional API key protection (set DISCORD_PROXY_KEY on server to enable)
    expectedKey := os.Getenv("DISCORD_PROXY_KEY")
    if expectedKey != "" {
        if r.Header.Get("X-API-KEY") != expectedKey {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
    }

    // Read raw request body and forward it unchanged so JSON structure is preserved.
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    if len(bodyBytes) == 0 {
        http.Error(w, "Empty request body", http.StatusBadRequest)
        return
    }

    // Determine content type once for validation AND for forwarding
    contentType := r.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/json"
    }

    // If JSON, do a light validation to avoid forwarding an empty Discord message
    if strings.Contains(contentType, "application/json") {
        var payload map[string]interface{}
        if err := json.Unmarshal(bodyBytes, &payload); err == nil {
            // Check for a non-empty `content` string or non-empty `embeds` array
            hasContent := false
            if c, ok := payload["content"].(string); ok {
                if strings.TrimSpace(c) != "" {
                    hasContent = true
                }
            }
            hasEmbeds := false
            if e, ok := payload["embeds"]; ok {
                if arr, ok := e.([]interface{}); ok && len(arr) > 0 {
                    hasEmbeds = true
                }
            }
            if !hasContent && !hasEmbeds {
                log.Printf("discord webhook rejected empty message payload: %s", string(bodyBytes))
                if os.Getenv("DEBUG_DISCORD_PROXY") == "true" {
                    http.Error(w, fmt.Sprintf("Rejected empty message payload: %s", string(bodyBytes)), http.StatusBadRequest)
                    return
                }
                http.Error(w, "Rejected empty message payload", http.StatusBadRequest)
                return
            }
        }
    }

    webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
    if webhookURL == "" {
        http.Error(w, "Webhook not configured", http.StatusInternalServerError)
        return
    }

    resp, err := http.Post(webhookURL, contentType, bytes.NewBuffer(bodyBytes))
    if err != nil {
        // network / DNS / TLS error
        log.Printf("discord webhook post error: %v", err)
        http.Error(w, "Failed to send webhook (post error)", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)
    respText := string(respBody)
    if len(respText) > 1000 {
        respText = respText[:1000] + "...(truncated)"
    }

    if resp.StatusCode >= 300 {
        // log full details server-side for diagnosis
        log.Printf("discord webhook failed: status=%d response=%s", resp.StatusCode, respText)

        // Optionally return more detail to the client when debugging is enabled
        if os.Getenv("DEBUG_DISCORD_PROXY") == "true" {
            http.Error(w, fmt.Sprintf("Failed to send webhook: status=%d body=%s", resp.StatusCode, respText), http.StatusBadGateway)
            return
        }

        http.Error(w, fmt.Sprintf("Failed to send webhook: status=%d", resp.StatusCode), http.StatusBadGateway)
        return
    }

    // success
    log.Printf("discord webhook forwarded successfully: status=%d", resp.StatusCode)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Webhook sent"))
}
