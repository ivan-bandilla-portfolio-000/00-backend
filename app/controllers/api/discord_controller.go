package api_controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
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

    webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
    if webhookURL == "" {
        http.Error(w, "Webhook not configured", http.StatusInternalServerError)
        return
    }

    contentType := r.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/json"
    }

    resp, err := http.Post(webhookURL, contentType, bytes.NewBuffer(bodyBytes))
    if err != nil {
        http.Error(w, "Failed to send webhook", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        http.Error(w, "Failed to send webhook", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Webhook sent"))
}
