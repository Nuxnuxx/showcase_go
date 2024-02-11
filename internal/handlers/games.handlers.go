package handlers

import (
	"net/http"
	"strconv"

	"github.com/Nuxnuxx/showcase_go/internal/services"
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
	page, err := strconv.Atoi(c.QueryParam("page"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page")
	}

	games, err := gh.GamesServices.GetGamesByPage(page)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(200, games)
}
