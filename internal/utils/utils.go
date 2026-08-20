package utils

import (
	"fmt"
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

// gets correct token url
// can be either prod or sandbox url
func TokenURL() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file for tokenURL")
	}

	tokenURL := os.Getenv("TOKEN_URL")

	if tokenURL == "" {
		panic("Cannot obtain tokenURL from .env file")
	}

	return tokenURL
}

// generic gets env variable and panics if not existant
func PanicIfNotExist(key string) string {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("could not obtain env variable from key %s", key))
	}

	envVal := os.Getenv(key)

	if envVal == "" {
		panic("Cannot obtain tokenURL from .env file")
	}

	return envVal
}
