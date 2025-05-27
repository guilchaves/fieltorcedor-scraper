package service

import (
	"fmt"

	"github.com/guilchaves/fieltorcedorbot/internal/core/domain"
	"github.com/guilchaves/fieltorcedorbot/internal/core/ports"
)

type GameService struct {
	gameRepo           ports.GameRepository
	notificationSender ports.NotificationSender
}

func NewGameService(
	gameRepo ports.GameRepository,
	notificationSender ports.NotificationSender,
) *GameService {
	return &GameService{
		gameRepo:           gameRepo,
		notificationSender: notificationSender,
	}
}

func (s *GameService) CheckForNewGames() ([]domain.Game, error) {
	fmt.Println("Verificando novos jogos...")
	games, err := s.gameRepo.FetchGames()
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar jogos: %w", err)
	}

	if len(games) == 0 {
		fmt.Println("Nenhum jogo encontrado.")
		return nil, nil
	}

	fmt.Printf("Encontrados %d jogos.Notificando...\n", len(games))
	err = s.NotifyAboutGames(games)
	if err != nil {
		return nil, fmt.Errorf("falha ao notificar jogos: %w", err)
	}
	return games, nil
}

func (s *GameService) NotifyAboutGames(games []domain.Game) error {
	for _, game := range games {
		err := s.notificationSender.SendNotification(game)
		if err != nil {
			return fmt.Errorf(
				"falha ao enviar notificação para o jogo %s x %s: %w",
				game.HomeTeam,
				game.AwayTeam,
				err,
			)
		}
		fmt.Printf("Notificação enviada para o jogo %s x %s\n", game.HomeTeam, game.AwayTeam)
	}
	return nil
}
