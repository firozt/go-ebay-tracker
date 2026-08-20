package main

import (
	"fmt"
	"go-ebay-tracker/internal/ebay"
	"os"
	"time"

	"github.com/joho/godotenv"
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

func main() {
	fmt.Println("================== running go-ebay-tracker ===================")

	// loading env vars

	fmt.Println("[LOGS] Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	APP_ID := os.Getenv("APP_ID")
	CERT_ID := os.Getenv("CERT_ID")
	if APP_ID == "" || CERT_ID == "" {
		panic("env variables not loaded for APP_ID or CERT_ID")
	}
	fmt.Println("[LOGS] Environemnet variables loaded!")

	// getting token
	fmt.Println("[LOGS] Obtaining refresh tokens...")
	tm := ebay.OAuthEbay(APP_ID, CERT_ID)
	token := tm.GetToken()
	fmt.Println("[LOGS] Tokens recieved!")

	searchLoop(token, 1)

	fmt.Println("================== stopping go-ebay-tracker ==================")
}

func searchLoop(token string, interval int) {
	iteration := 1
	for {
		fmt.Printf("========== ITERATION %d ==========\n", iteration)
		fmt.Println("[LOGS] Making search request")
		results, err := ebay.SearchMusicMagpie2DS(token)
		if err != nil {
			panic(err)
		}
		fmt.Printf("found %d listings\n", results.Total)
		for _, item := range results.ItemSummaries {
			fmt.Printf("- %s | %s %s | %s\n", item.Title, item.Price.Value, item.Price.Currency, item.ItemWebURL)
		}
		iteration++
		fmt.Println("==================================")
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
