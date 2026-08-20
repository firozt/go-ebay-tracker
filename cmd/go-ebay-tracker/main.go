package main

import (
	"fmt"
	"go-ebay-tracker/internal/ebay"
	"go-ebay-tracker/internal/utils"
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
	fmt.Println("Initialsiing application")

	// loading env vars
	fmt.Println("[LOGS] Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	APP_ID := utils.PanicIfNotExist("APP_ID")
	CERT_ID := os.Getenv("CERT_ID")
	fmt.Println("[LOGS] Environemnet variables loaded!")

	// getting initial token
	fmt.Println("[LOGS] Obtaining initial tokens...")
	tm := ebay.OAuthEbay(APP_ID, CERT_ID)
	fmt.Println("[LOGS] Tokens recieved!")
	fmt.Println("Initialsiing completed successfully")

	err = searchLoop(tm, 1)
	if err != nil {
		fmt.Printf("error within the search loop:\n%+v\n", err)
	}


	fmt.Println("================== stopping go-ebay-tracker ==================")
}

func searchLoop(tm *ebay.TokenManager, interval int) error {
	iteration := 1
	startTime := time.Now()
	for {
		fmt.Printf("========================= ITERATION %d =========================\n", iteration)
		token := tm.GetToken()

		fmt.Println("[LOGS] Making search request...")
		results, err := ebay.SearchMusicMagpie2DS(token)
		if err != nil {
			return err
		}
		fmt.Printf("found %d listings\n", results.Total)
		for _, item := range results.ItemSummaries {
			fmt.Printf("- %s | %s %s | %s\n", item.Title, item.Price.Value, item.Price.Currency, item.ItemWebURL)
		}
		fmt.Printf("Total elapsed time - %s\n", utils.ElapsedTime(startTime))
		iteration++
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
