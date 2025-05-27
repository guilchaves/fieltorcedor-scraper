package fieltorcedor

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/guilchaves/fieltorcedorbot/internal/core/domain"
)

type FielTorcedorScraper struct {
	baseURL string
}

func NewFielTorcedorScraper(baseURL string) *FielTorcedorScraper {
	return &FielTorcedorScraper{baseURL: baseURL}
}

func (s *FielTorcedorScraper) FetchGames() ([]domain.Game, error) {
	fmt.Println("Fetching games from Fiel Torcedor...")
	resp, err := http.Get(s.baseURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao fazer a requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code diferente de 200: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao parsear o HTML: %w", err)
	}

	games := []domain.Game{}

	doc.Find(".main-games article.tab_event").Each(func(i int, article *goquery.Selection) {
		game := domain.Game{}

		title := article.Find("h2.font24").Text()
		if strings.Contains(title, " X ") {
			times := strings.Split(title, " X ")
			if len(times) == 2 {
				game.HomeTeam = strings.TrimSpace(times[0])
				away := strings.TrimSpace(times[1])
				if idx := strings.Index(away, "("); idx != -1 {
					away = strings.TrimSpace(away[:idx])
				}
				game.AwayTeam = away
			}
		}

		game.Competition = strings.TrimSpace(
			article.Find("p.font16.font-gotham-medium.mb-0").First().Text(),
		)
		game.Round = strings.TrimSpace(article.Find("p.font16.font-gotham-medium").Eq(1).Text())
		game.Stadium = strings.TrimSpace(article.Find("p.font16.font-gotham-medium").Last().Text())

		dateTimeStr := strings.TrimSpace(article.Find("p.font20.font-gotham-medium.mb-0").Text())
		dateTimeStr = strings.ReplaceAll(dateTimeStr, "às", "")
		dateTimeStr = strings.TrimSpace(dateTimeStr)
		gameTime, err := time.Parse("02/01/2006 15:04", dateTimeStr)
		if err == nil {
			game.Date = gameTime
		}

		article.Find("span[data-balloon-pos='down']").
			Each(func(i int, selection *goquery.Selection) {
				categoryName := selection.Text()
				availabilityText := selection.AttrOr("aria-label", "")

				startTime, endTime, err := parseAvailabilityTimes(availabilityText)
				if err != nil {
					fmt.Printf(
						"Erro ao parsear os tempos de disponibilidade para a categoria %s: %v\n",
						categoryName,
						err,
					)
					return
				}

				category := domain.CategoryAvailability{
					Category:    domain.Category(categoryName),
					StartTime:   startTime,
					EndTime:     endTime,
					IsAvailable: true,
				}
				game.Categories = append(game.Categories, category)
			})

		game.ID = generateGameID(game)

		games = append(games, game)
	})

	return games, nil
}

func parseAvailabilityTimes(availabilityText string) (time.Time, time.Time, error) {
	parts := strings.Split(availabilityText, " até ")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf(
			"formato de texto de disponibilidade inválido: %s",
			availabilityText,
		)
	}

	startTimeStr := strings.ReplaceAll(parts[0], "De ", "")
	endTimeStr := parts[1]

	startTime, err := parseTime(startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("erro ao parsear o tempo de início: %w", err)
	}

	endTime, err := parseTime(endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("erro ao parsear o tempo de fim: %w", err)
	}

	return startTime, endTime, nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("02/01/2006 15:04", timeStr)
}

func generateGameID(game domain.Game) string {
	h := fnv.New64()
	_, err := h.Write([]byte(game.HomeTeam + game.AwayTeam + game.Date.String()))
	if err != nil {
		fmt.Printf("Erro ao gerar o hash: %v\n", err)
		return "0"
	}
	return fmt.Sprintf("%x", h.Sum64())
}
