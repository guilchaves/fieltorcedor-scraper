package notifiedgames

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/guilchaves/fieltorcedorbot/internal/core/ports"
	_ "github.com/lib/pq"
)

type SupabaseNotifiedGamesRepository struct {
	db *sql.DB
}

func NewSupabaseNotifiedGamesRepository(connStr string) (ports.NotifiedGamesRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Supabase: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("erro ao pingar o Supabase: %w", err)
	}

	log.Println("Conectado ao Supabase com sucesso!")
	return &SupabaseNotifiedGamesRepository{db: db}, nil
}

func (r *SupabaseNotifiedGamesRepository) IsGameNotified(gameID string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM notified_games WHERE game_id = $1", gameID).
		Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar se o jogo foi notificado: %w", err)
	}
	return count > 0, nil
}

func (r *SupabaseNotifiedGamesRepository) SaveNotifiedGame(gameID string) error {
	_, err := r.db.Exec(
		"INSERT INTO notified_games (game_id) VALUES ($1) ON CONFLICT (game_id) DO NOTHING",
		gameID,
	)
	if err != nil {
		return fmt.Errorf("erro ao salvar o jogo notificado: %w", err)
	}
	return nil
}

func (r *SupabaseNotifiedGamesRepository) Close() error {
	return r.db.Close()
}
