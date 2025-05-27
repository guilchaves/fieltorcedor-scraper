package service

import (
	"fmt"

	"github.com/guilchaves/fieltorcedorbot/internal/core/domain"
	"github.com/guilchaves/fieltorcedorbot/internal/core/ports"
)

type GameService struct {
	gameRepo            ports.GameRepository
	notificationSender  ports.NotificationSender
	notifiedGamesRepo   ports.NotifiedGamesRepository
}

func NewGameService(
	gameRepo ports.GameRepository,
	notificationSender ports.NotificationSender,
	notifiedGamesRepo ports.NotifiedGamesRepository, 
) *GameService {
	return &GameService{
		gameRepo:            gameRepo,
		notificationSender:  notificationSender,
		notifiedGamesRepo:   notifiedGamesRepo,
	}
}

func (s *GameService) GetAllGames() ([]domain.Game, error) {
	return s.gameRepo.FetchGames()
}

func (s *GameService) CheckForNewGames() ([]domain.Game, error) {
	fmt.Println("Verificando novos jogos...")
	games, err := s.gameRepo.FetchGames()
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar jogos: %w", err)
	}

	var newGames []domain.Game
	for _, game := range games {
		notified, err := s.notifiedGamesRepo.IsGameNotified(game.ID)
		if err != nil {
			fmt.Printf("Erro ao verificar se o jogo %s já foi notificado: %v\n", game.ID, err)
			continue
		}
		if !notified {
			newGames = append(newGames, game)
		}
	}

	if len(newGames) == 0 {
		fmt.Println("Nenhum novo jogo encontrado.")
		return nil, nil
	}

	fmt.Printf("Encontrados %d novos jogos. Notificando...\n", len(newGames))
	err = s.NotifyAboutGames(newGames)
	if err != nil {
		return nil, fmt.Errorf("falha ao notificar jogos: %w", err)
	}

	for _, game := range newGames {
		err := s.notifiedGamesRepo.SaveNotifiedGame(game.ID)
		if err != nil {
			fmt.Printf("Erro ao salvar o jogo notificado %s: %v\n", game.ID, err)
		}
	}

	return newGames, nil
}

func (s *GameService) NotifyAboutGames(games []domain.Game) error {
	for _, game := range games {
		err := s.notificationSender.SendNotification(game)
		if err != nil {
			return fmt.Errorf("falha ao enviar notificação para o jogo %s: %w", game.AwayTeam, err)
		}
	}
	return nil
}
