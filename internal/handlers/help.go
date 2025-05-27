package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HelpHandler(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Comandos disponíveis:\n/start - Mensagem de boas-vindas\n/jogos - Lista os próximos jogos do Corinthians\n/help - Mostra esta mensagem")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem de ajuda: %v", err)
	}
}
