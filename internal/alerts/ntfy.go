package alerts

import (
	"net/http"
	"strings"
)

// requires NTFY url, title and message
func SendNtfyAlert(NTFYUrl string, title string, message string) error {
	req, err := http.NewRequest("POST",NTFYUrl, strings.NewReader(message))
	if err != nil {
		return err
	}
	req.Header.Set("Title", title)
	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "rotating_light")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
