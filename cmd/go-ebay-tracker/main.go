package main

import (
	"fmt"
	"go-ebay-tracker/internal/ebay"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("================== running go-ebay-tracker ===================")
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	APP_ID := os.Getenv("APP_ID")
	CERT_ID := os.Getenv("CERT_ID")
	if APP_ID == "" || CERT_ID == "" {
		panic("env variables not loaded for APP_ID or CERT_ID")
	}
	tm := ebayapi.OAuthEbay(APP_ID, CERT_ID)
	fmt.Printf("token manager:\n%+v", tm)
	fmt.Println("================== stopping go-ebay-tracker ==================")
}
