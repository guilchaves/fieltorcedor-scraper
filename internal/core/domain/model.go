package domain

import "time"

type Category string

const (
	CategoryMinhaVida     Category = "Minha Vida"
	CategoryMinhaHistoria Category = "Minha História"
	CategoryMeuAmor       Category = "Meu Amor"
	CategoryMinhaCadeira  Category = "Minha Cadeira"
	CategoryNaoPagantes   Category = "Não Pagantes"
)

type Game struct {
	ID          string
	HomeTeam    string
	AwayTeam    string
	Competition string
	Round       string
	Date        time.Time
	Stadium     string
	Categories  []CategoryAvailability
}

type CategoryAvailability struct {
	Category    Category
	StartTime   time.Time
	EndTime     time.Time
	IsAvailable bool
}
