package main

import (
	"log"
	"os"
	"strconv"

	"github.com/guilchaves/fieltorcedorbot/internal/adapters/in"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/fieltorcedor"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notification"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notifiedgames"
	"github.com/guilchaves/fieltorcedorbot/internal/core/service"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	go func() {
		bot, err := tgbotapi.NewBotAPI(telegramToken)
		if err != nil {
			log.Fatalf("Erro ao iniciar o bot do Telegram: %v", err)
		}
		bot.Debug = false

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatalf("Erro ao obter updates do Telegram: %v", err)
		}

		for update := range updates {
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Olá! Buscando jogos do Corinthians para você...",
					)
					bot.Send(msg)

					games, err := gameService.GetAllGames()
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro ao buscar jogos.")
						bot.Send(msg)
						continue
					}
					if len(games) == 0 {
						msg := tgbotapi.NewMessage(
							update.Message.Chat.ID,
							"Nenhum jogo encontrado.",
						)
						bot.Send(msg)
						continue
					}

					for _, game := range games {
						err := telegram.SendNotificationToChat(
							game,
							strconv.FormatInt(update.Message.Chat.ID, 10),
						)
						if err != nil {
							log.Printf("Erro ao enviar notificação: %v", err)
						}
					}
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Comando não reconhecido.")
					bot.Send(msg)
				}
			}
		}
	}()

	select {}
}
