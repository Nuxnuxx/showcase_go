package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nuxnuxx/showcase_go/internal/services"
	"github.com/Nuxnuxx/showcase_go/internal/views/errors_pages"
	gamesviews "github.com/Nuxnuxx/showcase_go/internal/views/games_views"
	"github.com/labstack/echo/v4"
)

type GamesServices interface {
	GetGamesByPage(page int) ([]services.Game, error)
	GetGamesByID(id int) (services.GameFullDetail, error)
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
		c.Response().WriteHeader(http.StatusInternalServerError)
		return renderView(c, errors_pages.Error500Index())
	}

	if pageInt > 0 {
		return renderView(c, gamesviews.GamesList(games, pageInt))
	}

	return renderView(c, gamesviews.GameIndex(games, pageInt))
}

func (gh *GamesHandler) GetGameById(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return renderView(c, errors_pages.Error400Index())
	}

	game, err := gh.GamesServices.GetGamesByID(idInt)

	if err != nil {
		fmt.Println(err)
		c.Response().WriteHeader(http.StatusInternalServerError)
		return renderView(c, errors_pages.Error500Index())
	}

	return renderView(c, gamesviews.GamePageIndex(game))
}
