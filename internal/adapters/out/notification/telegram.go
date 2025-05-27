package notification

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/guilchaves/fieltorcedorbot/internal/core/domain"
)

type TelegramSender struct {
	botToken string
	chatID   string
}

func NewTelegramSender(botToken string, chatID string) *TelegramSender {
	return &TelegramSender{botToken: botToken, chatID: chatID}
}

func (s *TelegramSender) SendNotification(game domain.Game) error {
	message := fmt.Sprintf(
		"Novo jogo do Corinthians!\n\n"+
			"Adversário: %s\n"+
			"Campeonato: %s\n"+
			"Rodada: %s\n"+
			"Data: %s\n"+
			"Estádio: %s\n\n"+
			"Categorias Disponíveis:\n",
		game.AwayTeam,
		game.Competition,
		game.Round,
		game.Date.Format("02/01/2006 15:04"),
		game.Stadium,
	)

	for _, category := range game.Categories {
		message += fmt.Sprintf(
			"- %s (De %s até %s)\n",
			category.Category,
			category.StartTime.Format("02/01/2006 15:04"),
			category.EndTime.Format("02/01/2006 15:04"),
		)
	}

	err := s.sendMessage(message)
	if err != nil {
		return fmt.Errorf("falha ao enviar mensagem para o Telegram: %w", err)
	}

	return nil
}

func (s *TelegramSender) sendMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	encodedMessage := url.QueryEscape(message)
	apiURL += fmt.Sprintf("?chat_id=%s&text=%s", s.chatID, encodedMessage)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("falha ao fazer a requisição para a API do Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"status code diferente de 200 ao enviar mensagem para o Telegram: %d",
			resp.StatusCode,
		)
	}

	return nil
}

func GetTelegramBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func GetTelegramChatID() string {
	return os.Getenv("TELEGRAM_CHAT_ID")
}
