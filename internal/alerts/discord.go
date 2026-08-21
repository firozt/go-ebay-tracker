package alerts

import (
	"fmt"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// requires channel webhook url
func SendDiscordMessage(webhookURL, message string) error {
	log.Printf("[LOGS] Sending alert to discord\n")

	payload := WebhookPayload{Content: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord webhook error: status %d", resp.StatusCode)
	}
	return nil
}
