package main

import (
	"go-ebay-tracker/internal/ebay"
	"go-ebay-tracker/internal/utils"
	"io"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func setupLogger() (*os.File, error) {
	logFile, err := os.OpenFile("get.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, nil
}

func main() {
	log.Println("================== running go-ebay-tracker ===================")
	log.Println("Initialsiing application")

	log.Println("[LOGS] Setting up logger (current working dir -> go-ebay-tracker.log)")
	logFile, err := setupLogger()
	if err != nil {
		panic("Error unable to setup logger")
	}
	defer logFile.Close()
	log.Println("[LOGS] Logger setup successfully!")

	// loading env vars
	log.Println("[LOGS] Loading environment variables...")
	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	APP_ID := utils.PanicIfNotExist("APP_ID")
	CERT_ID := os.Getenv("CERT_ID")
	log.Println("[LOGS] Environemnet variables loaded!")

	// getting initial token
	log.Println("[LOGS] Obtaining initial tokens...")
	tm := ebay.OAuthEbay(APP_ID, CERT_ID)
	log.Println("[LOGS] Tokens recieved!")
	log.Println("Initialsiing completed successfully")

	err = searchLoop(tm, 1)
	if err != nil {
		log.Printf("error within the search loop:\n%+v\n", err)
	}

	log.Println("================== stopping go-ebay-tracker ==================")
}

func searchLoop(tm *ebay.TokenManager, interval int) error {
	iteration := 1
	startTime := time.Now()
	for {
		log.Printf("========================= ITERATION %d =========================\n", iteration)
		token := tm.GetToken()

		log.Println("[LOGS] Making search request...")
		results, err := ebay.SearchMusicMagpie2DS(token)
		if err != nil {
			return err
		}
		log.Printf("found %d listings\n", results.Total)
		for _, item := range results.ItemSummaries {
			log.Printf("- %s | %s %s | %s\n", item.Title, item.Price.Value, item.Price.Currency, item.ItemWebURL)
		}
		log.Printf("Total elapsed time - %s\n", utils.ElapsedTime(startTime))
		iteration++
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
