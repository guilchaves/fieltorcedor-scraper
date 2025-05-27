package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartHandler(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(
		chatID,
		"Olá! Bem-vindo ao Fiel Torcedor Bot.\nUse /jogos para ver os próximos jogos do Corinthians.",
	)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem de boas-vindas: %v", err)
	}
}
