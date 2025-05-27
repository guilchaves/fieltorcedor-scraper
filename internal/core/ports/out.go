package ports

import "github.com/guilchaves/fieltorcedorbot/internal/core/domain"

type GameRepository interface {
	FetchGames() ([]domain.Game, error)
}

type NotificationSender interface {
	SendNotification(game domain.Game) error
}
