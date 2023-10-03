package main

import (
	"github.com/joho/godotenv"
	"gitlab.com/whoophy/privy/config"
	"log"
)

func main() {
	// Load the .env file in the current directory
	godotenv.Load()

	err := config.NewConfig()
	if err != nil {
		log.Printf("Error while running program %s", err)
	}
}
