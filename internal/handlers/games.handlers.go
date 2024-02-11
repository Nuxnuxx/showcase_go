package handlers

import (
	"fmt"
	"strconv"

	"github.com/Nuxnuxx/showcase_go/internal/services"
	"github.com/Nuxnuxx/showcase_go/internal/views/errors_pages"
	gamesviews "github.com/Nuxnuxx/showcase_go/internal/views/games_views"
	"github.com/labstack/echo/v4"
)

type GamesServices interface {
	GetGamesByPage(page int) ([]services.Game, error)
}

func NewGamesHandlers(gs GamesServices) *GamesHandler {

	return &GamesHandler{
		GamesServices: gs,
	}
}

type GamesHandler struct {
	GamesServices GamesServices
}

func (gh *GamesHandler) GetGamesByPage(c echo.Context) error {
	page := c.QueryParam("page")

	if page == "" {
		page = "0"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return renderView(c, errors_pages.Error400Index())
	}

	games, err := gh.GamesServices.GetGamesByPage(pageInt)

	if err != nil {
		return renderView(c, errors_pages.Error500Index())
	}

	fmt.Println(pageInt)

	if pageInt > 0 {
			return renderView(c, gamesviews.GamesList(games, pageInt))
	}
	
	return renderView(c, gamesviews.GameIndex(games, pageInt))
}
