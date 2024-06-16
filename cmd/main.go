package main

import (
	"log"
	"os"

	"projectik/internal/database"
	"projectik/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := telegram.NewBot(botToken, db)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
