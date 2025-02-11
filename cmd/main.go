package main

import (
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"tgbot-url-keeper/internal/repository/storage"
	"tgbot-url-keeper/internal/telegram"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}

	if err := storage.Init(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Токен бота не указан")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	telegram.SetupBot(bot)

	log.Printf("Бот запущен. токен: %s", os.Getenv("TELEGRAM_BOT_TOKEN"))
	bot.Start()

}
