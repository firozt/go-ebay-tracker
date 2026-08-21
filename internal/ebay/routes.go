package ebay

import (
	"encoding/json"
	"fmt"
	"go-ebay-tracker/internal/alerts"
	"go-ebay-tracker/internal/utils"
	"io"
	"net/http"
	"net/url"
)

type SearchResponse struct {
	Total         int           `json:"total"`
	ItemSummaries []ItemSummary `json:"itemSummaries"`
}

type ItemSummary struct {
	ItemID     string `json:"itemId"`
	Title      string `json:"title"`
	Price      Price  `json:"price"`
	ItemWebURL string `json:"itemWebUrl"`
}

type Price struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

func SearchMusicMagpie2DS(token string) (*SearchResponse, error) {
	// build up request
	var baseURL string
	// sandbox market or real ebay market
	if utils.IsDevMode() {
		baseURL = "https://api.sandbox.ebay.com/buy/browse/v1/item_summary/search"
	} else {
		baseURL = "https://api.ebay.com/buy/browse/v1/item_summary/search"
	}

	params := url.Values{}
	params.Set("q", "2ds xl")
	// params.Set("category_ids", "139971")
	params.Set("filter", "sellers:{musicmagpie}")
	params.Set("sort", "newlyListed")
	// params.Set("sort", "price")
	params.Set("limit", "5")

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-EBAY-C-MARKETPLACE-ID", "EBAY_GB")

	// send+handle request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ebay api error %d: %s", resp.StatusCode, body)
	}

	// marshal body as searchResponse object
	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// send notifs (discord + Ntfy)
	if len(result.ItemSummaries) > 0 {
		discordWebhook := utils.PanicIfNotExist("DISCORD_WEBHOOK")
		firstItemURL := result.ItemSummaries[0].ItemWebURL
		alerts.SendDiscordMessage(discordWebhook, fmt.Sprintf("New listing for 2ds xl from musicmagpie @everyone %s", firstItemURL))

		NTFYUrl := utils.PanicIfNotExist("NTFY_URL")
		alerts.SendNtfyAlert(
			NTFYUrl,
			"New 2DS Listing!",
			fmt.Sprintf("New 2ds xl listing for musicmagpie! -> %s", firstItemURL),
		)
	}

	return &result, nil
}

