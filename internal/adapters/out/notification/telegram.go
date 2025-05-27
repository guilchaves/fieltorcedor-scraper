package notification

import (
	"fmt"
	"io"
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
		"🦅 *Novo jogo do Corinthians\\!* 🦅\n\n"+
			"*Adversário:* %s\n"+
			"*Campeonato:* %s\n"+
			"*Rodada:* %s\n"+
			"*Data:* %s\n"+
			"*Estádio:* %s\n\n"+
			"*Categorias Disponíveis:*\n",
		escapeMarkdown(game.AwayTeam),
		escapeMarkdown(game.Competition),
		escapeMarkdown(game.Round),
		escapeMarkdown(game.Date.Format("02/01/2006 15:04")),
		escapeMarkdown(game.Stadium),
	)

	for _, category := range game.Categories {
		message += "\\- " + escapeMarkdown(string(category.Category)) +
			" \\(De " + escapeMarkdown(category.StartTime.Format("02/01/2006 15:04")) +
			" até " + escapeMarkdown(category.EndTime.Format("02/01/2006 15:04")) + "\\)\n"
	}

	err := s.sendMessage(message)
	if err != nil {
		return fmt.Errorf("falha ao enviar mensagem para o Telegram: %w", err)
	}

	return nil
}

func (s *TelegramSender) sendMessage(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	data := url.Values{}
	data.Set("chat_id", s.chatID)
	data.Set("text", message)
	data.Set("parse_mode", "MarkdownV2")

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("falha ao fazer a requisição para a API do Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"status code diferente de 200 ao enviar mensagem para o Telegram: %d\nCorpo da resposta: %s",
			resp.StatusCode,
			string(body),
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

func (s *TelegramSender) SendNotificationToChat(game domain.Game, chatID string) error {
	message := fmt.Sprintf(
		"🦅 *Novo jogo do Corinthians\\!* 🦅\n\n"+
			"*Adversário:* %s\n"+
			"*Campeonato:* %s\n"+
			"*Rodada:* %s\n"+
			"*Data:* %s\n"+
			"*Estádio:* %s\n\n",
		escapeMarkdown(game.AwayTeam),
		escapeMarkdown(game.Competition),
		escapeMarkdown(game.Round),
		escapeMarkdown(game.Date.Format("02/01/2006 15:04")),
		escapeMarkdown(game.Stadium),
	)

	if len(game.Categories) > 0 {
		message += "*Categorias Disponíveis:*\n"
		for _, category := range game.Categories {
			message += fmt.Sprintf(
				"\\- %s \\(De %s até %s\\)\n",
				escapeMarkdown(string(category.Category)),
				escapeMarkdown(category.StartTime.Format("02/01/2006 15:04")),
				escapeMarkdown(category.EndTime.Format("02/01/2006 15:04")),
			)
		}
	}

	return s.sendMessageToChat(message, chatID)
}

func (s *TelegramSender) sendMessageToChat(message string, chatID string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "MarkdownV2")

	resp, err := http.PostForm(apiURL, data)
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
