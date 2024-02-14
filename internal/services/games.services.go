package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nuxnuxx/showcase_go/internal/database"
)

func NewGamesServices(g Game, gStore database.Store, apiKey string) *GameService {

	return &GameService{
		Game:      g,
		GameStore: gStore,
		ApiKey:    apiKey,
	}
}

type GameService struct {
	Game      Game
	GameStore database.Store
	ApiKey    string
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
	Rating           float32     `json:"rating"`
	RatingTop        int         `json:"rating_top"`
	Ratings          interface{} `json:"ratings"`
	RatingsCount     int         `json:"ratings_count"`
	ReviewsTextCount int         `json:"reviews_text_count"`
	Added            int         `json:"added"`
	AddedByStatus    interface{} `json:"added_by_status"`
	Metacritic       int         `json:"metacritic"`
	Playtime         int         `json:"playtime"`
	SuggestionsCount int         `json:"suggestions_count"`
	Updated          string      `json:"updated"`
	EsrbRating       EsrbRating  `json:"esrb_rating"`
	Platforms        []Platform  `json:"platforms"`
}

type MetacriticPlatform struct {
	Metascore int    `json:"metascore"`
	URL       string `json:"url"`
}

type GameFullDetail struct {
	ID                        int                    `json:"id"`
	Slug                      string                 `json:"slug"`
	Name                      string                 `json:"name"`
	NameOriginal              string                 `json:"name_original"`
	Description               string                 `json:"description"`
	Metacritic                int                    `json:"metacritic"`
	MetacriticPlatforms       []MetacriticPlatform   `json:"metacritic_platforms"`
	Released                  string                 `json:"released"`
	TBA                       bool                   `json:"tba"`
	Updated                   string                 `json:"updated"`
	BackgroundImage           string                 `json:"background_image"`
	BackgroundImageAdditional string                 `json:"background_image_additional"`
	Website                   string                 `json:"website"`
	Rating                    float32                `json:"rating"`
	RatingTop                 int                    `json:"rating_top"`
	Reactions                 map[string]interface{} `json:"reactions"`
	Added                     int                    `json:"added"`
	AddedByStatus             map[string]interface{} `json:"added_by_status"`
	Playtime                  int                    `json:"playtime"`
	ScreenshotsCount          int                    `json:"screenshots_count"`
	MoviesCount               int                    `json:"movies_count"`
	CreatorsCount             int                    `json:"creators_count"`
	AchievementsCount         int                    `json:"achievements_count"`
	ParentAchievementsCount   int                    `json:"parent_achievements_count"`
	RedditURL                 string                 `json:"reddit_url"`
	RedditName                string                 `json:"reddit_name"`
	RedditDescription         string                 `json:"reddit_description"`
	RedditLogo                string                 `json:"reddit_logo"`
	RedditCount               int                    `json:"reddit_count"`
	TwitchCount               int                    `json:"twitch_count"`
	YoutubeCount              int                    `json:"youtube_count"`
	ReviewsTextCount          int                    `json:"reviews_text_count"`
	RatingsCount              int                    `json:"ratings_count"`
	SuggestionsCount          int                    `json:"suggestions_count"`
	AlternativeNames          []string               `json:"alternative_names"`
	MetacriticURL             string                 `json:"metacritic_url"`
	ParentsCount              int                    `json:"parents_count"`
	AdditionsCount            int                    `json:"additions_count"`
	GameSeriesCount           int                    `json:"game_series_count"`
	ESRBRating                EsrbRating             `json:"esrb_rating"`
	Platforms                 []Platform             `json:"platforms"`
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
	if page > 0 {
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

func (gs *GameService) GetGamesByID(id int) (GameFullDetail, error) {
	// Make the url
	builder := strings.Builder{}
	builder.WriteString("https://api.rawg.io/api/games/")
	builder.WriteString(strconv.Itoa(id))
	builder.WriteString("?key=")
	builder.WriteString(gs.ApiKey)

	resp, err := http.Get(builder.String())

	if err != nil {
		return GameFullDetail{}, fmt.Errorf("Error making request: %v", err)
	}

	defer resp.Body.Close()

	// This part bind the response to the struct
	var response GameFullDetail

	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &response); err != nil {
		return GameFullDetail{}, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return response, nil
}

func (gs *GameService) LikeGameByID(id, idUser int) error {

	query := `INSERT INTO user_liked_games (user_id, liked_game_id) VALUES($1, $2)`

	_, err := gs.GameStore.Db.Exec(
		query,
		idUser,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (gs *GameService) GetGamesLikedByUser(id int) ([]GameFullDetail, error) {
	// Check if there are any games liked by the user first
	query := `SELECT liked_game_id FROM user_liked_games WHERE user_id = $1`

	rows, err := gs.GameStore.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameFullDetail

	// Check if there are no rows returned
	if !rows.Next() {
		// No games liked by the user, return an empty slice
		return games, nil
	}

	// Iterate over the result set
	for rows.Next() {
		var likedGameID int
		if err := rows.Scan(&likedGameID); err != nil {
			return nil, err
		}

		// Get game details by ID
		game, err := gs.GetGamesByID(likedGameID)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}
