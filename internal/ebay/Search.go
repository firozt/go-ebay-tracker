package ebay

import (
	"encoding/json"
	"fmt"
	"go-ebay-tracker/internal/alerts"
	"go-ebay-tracker/internal/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type SortType string

const (
	PriceAsc      SortType = "price"
	PriceDesc     SortType = "-price"
	Distance      SortType = "distance"
	NewlyListed   SortType = "newlyListed"
	EndingSoonest SortType = "endingSoonest"
	BestMatch     SortType = ""
)

type SearchFilter struct {
	Query                 string   `json:"q"`
	CategoryIDs           []string `json:"category_ids"` // csv
	EbayProductID         string   `json:"epid"`
	Filter                string   // TODO: figure out how to do this lol
	GlobalTradeItemNumber string   `json:"gtin"`
	CharityIDs            []string `json:"charity_ids"`
	Limit                 uint     `json:"limit"`  // max 1-200, default 50
	Offset                uint     `json:"offset"` // default 0, max 10_000
	Sort                  SortType `json:"sort_type"`
}

func (sf SearchFilter) ConstructSearchParams() (url.Values, error) {
	// at least one of these must be set per eBay's Browse API requirements
	if sf.Query == "" && len(sf.CategoryIDs) == 0 && sf.EbayProductID == "" &&
		sf.GlobalTradeItemNumber == "" && len(sf.CharityIDs) == 0 {
		return nil, fmt.Errorf("search requires at least one of: query, category_ids, epid, gtin, charity_ids")
	}

	if sf.Limit > 200 {
		log.Printf("[WARN] maximum value for limit in search request is 200, %d was given. setting value to cap of 200", sf.Limit)
		sf.Limit = 200
	}
	if sf.Offset > 10_000 {
		log.Printf("[WARN] maximum value for Offset in search request is 10_000, %d was given. setting value to cap of 10_000", sf.Offset)
		sf.Offset = 10_000
	}

	params := url.Values{}

	// build params, only defined ones
	if sf.Query != "" {
		params.Set("q", sf.Query)
	}
	if len(sf.CategoryIDs) > 0 {
		params.Set("category_ids", strings.Join(sf.CategoryIDs, ","))
	}
	if sf.EbayProductID != "" {
		params.Set("epid", sf.EbayProductID)
	}
	if sf.Filter != "" {
		params.Set("filter", sf.Filter)
	}
	if sf.GlobalTradeItemNumber != "" {
		params.Set("gtin", sf.GlobalTradeItemNumber)
	}
	if sf.Sort != "" {
		params.Set("sort", string(sf.Sort))
	}
	if len(sf.CharityIDs) > 0 {
		params.Set("charity_ids", strings.Join(sf.CharityIDs, ","))
	}
	if sf.Limit > 0 {
		params.Set("limit", strconv.FormatUint(uint64(sf.Limit), 10))
	}
	if sf.Offset > 0 {
		params.Set("offset", strconv.FormatUint(uint64(sf.Offset), 10))
	}

	return params, nil
}

func SearchListing(tm *TokenManager, queryParams SearchFilter) (*SearchResponse, error) {
	// build up request
	var baseURL string
	// sandbox market or real ebay market
	if utils.IsDevMode() {
		baseURL = "https://api.sandbox.ebay.com/buy/browse/v1/item_summary/search"
	} else {
		baseURL = "https://api.ebay.com/buy/browse/v1/item_summary/search"
	}

	params, err := queryParams.ConstructSearchParams()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tm.GetToken())
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

	// TODO refactor this to make it opt in
	// send notifs (discord + Ntfy)
	if len(result.ItemSummaries) > 0 {
		discordWebhook := utils.PanicIfNotExist("DISCORD_WEBHOOK")
		message := buildAlertMessage(queryParams, result.ItemSummaries)

		if err := alerts.SendDiscordMessage(discordWebhook,message); err != nil {
			log.Println("[ERROR] discord message failed to send:", err)
		}

		NTFYUrl := utils.PanicIfNotExist("NTFY_URL")
		if err := alerts.SendNtfyAlert(NTFYUrl, "New musicmagpie listing for 2ds xl", message,); err != nil {
			log.Println("[ERROR] ntfy alert message failed to send:", err)
		}

		if err := alerts.TriggerHomeMiniAlert(alerts.HomeMiniAlertRequest{}); err != nil {
			log.Println("[ERROR] home mini alert failed:", err)
		}
	}

	return &result, nil
}

// prints out a message like the following when multiple items returned:
// New listing for "2ds xl" — Nintendo 2DS XL Pokéball Edition (£89.99GBP) — https://www.ebay.co.uk/itm/123456789
// ...and 2 more
//
// otherwise:
// New listing for "2ds xl" — Nintendo 2DS XL Pokéball Edition (£89.99GBP) — https://www.ebay.co.uk/itm/123456789
func buildAlertMessage(sf SearchFilter, items []ItemSummary) string {
	desc := sf.Query
	if desc == "" {
		desc = "your search"
	}

	newest := items[0]
	msg := fmt.Sprintf("New listing for %q — %s (%s%s) — %s",
		desc, newest.Title, newest.Price.Value, newest.Price.Currency, newest.ItemWebURL)

	if len(items) > 1 {
		msg += fmt.Sprintf("\n...and %d more", len(items)-1)
	}

	return msg
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
		alerts.SendNtfyAlert(NTFYUrl, "New 2DS Listing!",
			fmt.Sprintf("New 2ds xl listing for musicmagpie! -> %s", firstItemURL),
		)

		if err := alerts.TriggerHomeMiniAlert(alerts.HomeMiniAlertRequest{}); err != nil {
			log.Println("[ERROR] home mini alert failed:", err)
		}
	}

	return &result, nil
}
