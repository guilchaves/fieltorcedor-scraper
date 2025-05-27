package handlers

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/guilchaves/fieltorcedorbot/internal/core/service"
	"github.com/guilchaves/fieltorcedorbot/internal/core/ports"
)

func JogosHandler(
	bot *tgbotapi.BotAPI,
	chatID int64,
	gameService *service.GameService,
	telegram ports.NotificationSender,
) {
	msg := tgbotapi.NewMessage(chatID, "Buscando jogos do Corinthians para você...")
	bot.Send(msg)

	games, err := gameService.GetAllGames()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Erro ao buscar jogos.")
		bot.Send(msg)
		return
	}
	if len(games) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Nenhum jogo encontrado.")
		bot.Send(msg)
		return
	}

	for _, game := range games {
		err := telegram.SendNotificationToChat(game, strconv.FormatInt(chatID, 10))
		if err != nil {
			log.Printf("Erro ao enviar notificação: %v", err)
		}
	}
}
