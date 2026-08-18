package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// checks env variable for devmode bool
// failsafe, if undefined or unexpected shape
// will result in true (use devmode sandbox options)
func IsDevMode() bool {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	if os.Getenv("DEVMODE") == "0" {
		return false
	}

	return true
}
