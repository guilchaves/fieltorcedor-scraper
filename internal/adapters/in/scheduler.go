package in

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/guilchaves/fieltorcedorbot/internal/core/ports"
)

type Scheduler struct {
	gameUseCase ports.GameUseCase
	cron        *gocron.Scheduler
	spec        string
	onceJobs    map[string]struct{}
	mu          sync.Mutex
}

func NewScheduler(gameUseCase ports.GameUseCase, spec string) *Scheduler {
	s := gocron.NewScheduler(time.Local)
	return &Scheduler{
		gameUseCase: gameUseCase,
		cron:        s,
		spec:        spec,
		onceJobs:    make(map[string]struct{}),
		mu:          sync.Mutex{},
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

func (s *Scheduler) ScheduleOnce(t time.Time, jobName string, job func()) error {
	if t.Before(time.Now()) {
		log.Printf("Horário %s já passou, ignorando agendamento para %s",
			t.Format("02/01/2006 15:04"), jobName)
		return nil
	}

	key := t.Format("200601021504") + "_" + jobName

	s.mu.Lock()
	if _, exists := s.onceJobs[key]; exists {
		s.mu.Unlock()
		log.Printf("Job %s já agendado para %s, ignorando",
			jobName, t.Format("02/01/2006 15:04"))
		return nil
	}

	s.onceJobs[key] = struct{}{}
	s.mu.Unlock()

	_, err := s.cron.At(t).Do(func() {
		job()

		s.mu.Lock()
		delete(s.onceJobs, key)
		s.mu.Unlock()

		log.Printf("Job %s executado e removido da lista", jobName)
	})

	if err != nil {
		return fmt.Errorf("erro ao agendar job %s: %w", jobName, err)
	}

	log.Printf("Job %s agendado com sucesso para %s",
		jobName, t.Format("02/01/2006 15:04"))
	return nil
}
