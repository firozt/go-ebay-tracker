package ebayapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// holds response type for
// oauth post request to obtain initial tokens
type EbayOAuthResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   time.Duration `json:"expires_in"`
}

// obtains tokens from ebay OAuth, documented in
// https://developer.ebay.com/develop/guides/sell/authorization
//
// returns token manager, from which token can be obtained safely
// function may panic given this fails as its crucial for application to work
func OAuthEbay(clientID string, clientSecret string) *TokenManager {
	secret := fmt.Sprintf(clientID, ":", clientSecret)
	secretBase64 := "Basic " + base64.StdEncoding.EncodeToString([]byte(secret))

	response, err := sendEbayOAuthPostRequest(secretBase64)

	if err != nil {
		panic("Failed during OAuth2 initialisation, could not obtain tokens")
	}

	TokenManager := NewTokenManager(secretBase64, response.AccessToken, response.ExpiresIn)

	return TokenManager
}

// resfresh token, skip calculation logic
func RefreshToken(b64secret string) (EbayOAuthResponse, error) {
	return sendEbayOAuthPostRequest(b64secret)
}

// request shape
// POST /identity/v1/oauth2/token
// Host: api.ebay.com
// Content-Type: application/x-www-form-urlencoded
// Authorization: Basic base64(client_id:client_secret)
//
// attatch secret base64 encoded
// grant_type=client_credentials&scope=https://api.ebay.com/oauth/api_scope
func sendEbayOAuthPostRequest(b64Secret string) (EbayOAuthResponse, error) {
	postURL := "https://api.sandbox.ebay.com/identity/v1/oauth2/token" // sandbox/dev
	// postURL := "https://api.ebay.com/identity/v1/oauth2/token" // prod

	body := strings.NewReader(url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"https://api.ebay.com/oauth/api_scope"},
	}.Encode())

	req, err := http.NewRequest(http.MethodPost, postURL, body)
	if err != nil {
		return EbayOAuthResponse{}, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+b64Secret)

	// 10s timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return EbayOAuthResponse{}, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return EbayOAuthResponse{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, respBody)
	}

	var result EbayOAuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return EbayOAuthResponse{}, fmt.Errorf("decoding response: %w", err)
	}

	return result, nil
}
