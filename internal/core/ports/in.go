package ports

import "github.com/guilchaves/fieltorcedorbot/internal/core/domain"

type GameUseCase interface {
	CheckForNewGames() ([]domain.Game, error)

	NotifyAboutGames(games []domain.Game) error
}
