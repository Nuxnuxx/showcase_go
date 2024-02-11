package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Nuxnuxx/showcase_go/internal/database"
)

func NewGamesServices(g Game, gStore database.Store, apiKey string) *GameService {

	return &GameService{
		Game:      g,
		GameStore: gStore,
		ApiKey: apiKey,
	}
}

type GameService struct {
	Game      Game
	GameStore database.Store
	ApiKey string
}

type EsrbRating struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Platform struct {
	ID           int    `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	ReleasedAt   string `json:"released_at"`
	Requirements struct {
		Minimum     string `json:"minimum"`
		Recommended string `json:"recommended"`
	} `json:"requirements"`
}

type Game struct {
	ID               int         `json:"id"`
	Slug             string      `json:"slug"`
	Name             string      `json:"name"`
	Released         string      `json:"released"`
	Tba              bool        `json:"tba"`
	BackgroundImage  string      `json:"background_image"`
	Rating           float32         `json:"rating"`
	RatingTop        int         `json:"rating_top"`
	Ratings          interface{} `json:"ratings"`
	RatingsCount     int         `json:"ratings_count"`
	ReviewsTextCount int      `json:"reviews_text_count"`
	Added            int         `json:"added"`
	AddedByStatus    interface{} `json:"added_by_status"`
	Metacritic       int         `json:"metacritic"`
	Playtime         int         `json:"playtime"`
	SuggestionsCount int         `json:"suggestions_count"`
	Updated          string   `json:"updated"`
	EsrbRating       EsrbRating  `json:"esrb_rating"`
	Platforms        []Platform  `json:"platforms"`
}

type Response struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Game `json:"results"`
}

func (gs *GameService) GetGamesByPage(page int) ([]Game, error) {
	// Make the url
	builder := strings.Builder{}
	builder.WriteString("https://api.rawg.io/api/games?key=")
	builder.WriteString(gs.ApiKey)

	// If page is not the first page, add the page number to the request
	if (page > 0) {
		builder.WriteString(fmt.Sprintf("&page=%d", page))
	}

	// Make the request
	resp, err := http.Get(builder.String())

	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}

	defer resp.Body.Close()

	// This part bind the response to the struct
	var response Response

	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return response.Results, nil
}
