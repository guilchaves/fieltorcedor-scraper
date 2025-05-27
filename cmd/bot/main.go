package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/in"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/fieltorcedor"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notification"
	"github.com/guilchaves/fieltorcedorbot/internal/core/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Erro ao carregar o arquivo .env. Variáveis de ambiente podem não estar definidas.")
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramToken == "" || telegramChatID == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN e TELEGRAM_CHAT_ID devem estar definidos nas variáveis de ambiente")
	}

	siteURL := "https://www.fieltorcedor.com.br/"

	scraper := fieltorcedor.NewFielTorcedorScraper(siteURL)
	telegram := notification.NewTelegramSender(telegramToken, telegramChatID)

	gameService := service.NewGameService(scraper, telegram)

	scheduler := in.NewScheduler(gameService, "0 7 * * *")
	err = scheduler.Start()
	if err != nil {
		log.Fatalf("Erro ao iniciar o scheduler: %v", err)
	}

	select {}
}
