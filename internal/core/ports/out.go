package ports

import "github.com/guilchaves/fieltorcedorbot/internal/core/domain"

type GameRepository interface {
	FetchGames() ([]domain.Game, error)
}

type NotificationSender interface {
	SendNotification(game domain.Game) error
	SendNotificationToChat(game domain.Game, chatID string) error
}

type NotifiedGamesRepository interface {
	SaveNotifiedGame(gameID string) error
	IsGameNotified(gameID string) (bool, error)
}
