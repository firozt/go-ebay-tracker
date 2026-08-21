package alerts
// this will is a general template for home mini support (webhook)

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HomeMiniAlertRequest struct {
	SoundURL   string  `json:"sound_url,omitempty"`
	Repeat     int     `json:"repeat,omitempty"`
	Volume     float64 `json:"volume,omitempty"`
	GapSeconds float64 `json:"gap_seconds,omitempty"`
}

func TriggerHomeMiniAlert(req HomeMiniAlertRequest) error {
	fmt.Printf("[LOGS] Sending alert to Home Mini (127.0.0.1:5005/alert)\n")
	req.Repeat = 15
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://127.0.0.1:5005/alert", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("home mini alert failed: status %d", resp.StatusCode)
	}
	return nil
}
