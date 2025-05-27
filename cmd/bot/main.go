package main

import (
	"log"
	"os"

	"github.com/guilchaves/fieltorcedorbot/internal/adapters/in"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/fieltorcedor"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notification"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notifiedgames"
	"github.com/guilchaves/fieltorcedorbot/internal/core/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(
			"Erro ao carregar o arquivo .env. Variáveis de ambiente podem não estar definidas.",
		)
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramToken == "" || telegramChatID == "" {
		log.Fatal(
			"TELEGRAM_BOT_TOKEN e TELEGRAM_CHAT_ID devem estar definidos nas variáveis de ambiente",
		)
	}

	siteURL := "https://www.fieltorcedor.com.br/"

	scraper := fieltorcedor.NewFielTorcedorScraper(siteURL)
	telegram := notification.NewTelegramSender(telegramToken, telegramChatID)

	notifiedRepo := notifiedgames.NewFileNotifiedGamesRepository("notified_games.txt")
	gameService := service.NewGameService(scraper, telegram, notifiedRepo)

	scheduler := in.NewScheduler(gameService, "* * * * *")
	err = scheduler.Start()
	if err != nil {
		log.Fatalf("Erro ao iniciar o scheduler: %v", err)
	}

	dailyScheduler := in.NewScheduler(gameService, "0 11 * * *")
	if err := dailyScheduler.Start(); err != nil {
		log.Fatalf("erro ao iniciar o scheduler diário: %v", err)
	}

	go func() {
		if _, err := gameService.CheckForNewGames(); err != nil {
			log.Printf("erro na verificação inicial: %v", err)
		}
	}()

	select {}
}
