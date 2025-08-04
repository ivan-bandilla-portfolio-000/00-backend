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
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	reqBody, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode >= 300 {
		http.Error(w, "Failed to send webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook sent"))
}
