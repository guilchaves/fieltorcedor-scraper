package in

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/guilchaves/fieltorcedorbot/internal/core/ports"
)

type Scheduler struct {
	gameUseCase ports.GameUseCase
	cron        *gocron.Scheduler
	spec        string
}

func NewScheduler(gameUseCase ports.GameUseCase, spec string) *Scheduler {
	s := gocron.NewScheduler(time.Local)
	return &Scheduler{
		gameUseCase: gameUseCase,
		cron:        s,
		spec:        spec,
	}
}

func (s *Scheduler) Start() error {
	_, err := s.cron.Cron(s.spec).Do(s.gameUseCase.CheckForNewGames)
	if err != nil {
		return fmt.Errorf("falha ao agendar a tarefa: %w", err)
	}

	s.cron.StartAsync()
	log.Println("Scheduler iniciado com o spec:", s.spec)
	return nil
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler interrompido")
}
