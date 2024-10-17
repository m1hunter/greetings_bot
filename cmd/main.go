package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"projectik/internal/database"
	"projectik/internal/notification"
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
	defer db.Close()

	if err := database.CreateTables(db); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := notification.NewBot(botToken, db)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
