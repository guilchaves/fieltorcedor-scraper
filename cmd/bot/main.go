package main

import (
	"log"
	"os"

	"github.com/guilchaves/fieltorcedorbot/internal/adapters/in"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/fieltorcedor"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notification"
	"github.com/guilchaves/fieltorcedorbot/internal/adapters/out/notifiedgames"
	"github.com/guilchaves/fieltorcedorbot/internal/core/service"
	"github.com/guilchaves/fieltorcedorbot/internal/handlers"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Erro ao carregar o arquivo .env. Variáveis de ambiente podem não estar definidas.")
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	supabaseConnString := os.Getenv("SUPABASE_CONNECTION_STRING") 
	if telegramToken == "" || telegramChatID == "" || supabaseConnString == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID e SUPABASE_CONNECTION_STRING devem estar definidos nas variáveis de ambiente")
	}

	siteURL := "https://www.fieltorcedor.com.br/"

	scraper := fieltorcedor.NewFielTorcedorScraper(siteURL)
	telegram := notification.NewTelegramSender(telegramToken, telegramChatID)

	notifiedRepo, err := notifiedgames.NewSupabaseNotifiedGamesRepository(supabaseConnString)
	if err != nil {
		log.Fatalf("Erro ao criar repositório Supabase: %v", err)
	}
	if closer, ok := notifiedRepo.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	gameService := service.NewGameService(scraper, telegram, notifiedRepo)

	schedules := []struct {
		name string
		spec string
	}{
		{"11h", "0 11 * * *"},
		{"15h", "0 15 * * *"},
		{"17h", "0 17 * * *"},
		{"18h", "0 18 * * *"},
	}

	for _, schedule := range schedules {
		scheduler := in.NewScheduler(gameService, schedule.spec)
		if err := scheduler.Start(); err != nil {
			log.Fatalf("erro ao iniciar o scheduler das %s: %v", schedule.name, err)
		}
		log.Printf("Scheduler das %s iniciado com sucesso", schedule.name)
	}

	go func() {
		log.Println("Realizando verificação inicial...")
		if _, err := gameService.CheckForNewGames(); err != nil {
			log.Printf("erro na verificação inicial: %v", err)
		} else {
			log.Println("Verificação inicial concluída com sucesso")
		}
	}()

	go func() {
		bot, err := tgbotapi.NewBotAPI(telegramToken)
		if err != nil {
			log.Fatalf("Erro ao iniciar o bot do Telegram: %v", err)
		}
		log.Printf("Bot do Telegram iniciado com sucesso: @%s", bot.Self.UserName)
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
				chatID := update.Message.Chat.ID
				switch update.Message.Command() {
				case "start":
					handlers.StartHandler(bot, chatID)
				case "jogos":
					handlers.JogosHandler(bot, chatID, gameService, telegram)
				case "help":
					handlers.HelpHandler(bot, chatID)
				default:
					msg := tgbotapi.NewMessage(chatID, "Comando não reconhecido.")
					bot.Send(msg)
				}
			}
		}
	}()

	log.Println("Bot iniciado e aguardando comandos...")
	select {}
}
