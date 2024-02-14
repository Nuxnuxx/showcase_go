package handlers

import (
	"net/http"
	"strconv"

	"github.com/Nuxnuxx/showcase_go/internal/services"
	"github.com/Nuxnuxx/showcase_go/internal/views/errors_pages"
	gamesviews "github.com/Nuxnuxx/showcase_go/internal/views/games_views"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewUserInteractionHandler(as AuthServices, gs GamesServices) *UserInteractionHandler {
	return &UserInteractionHandler{
		GamesServices: gs,
		AuthServices:  as,
	}
}

type UserInteractionHandler struct {
	GamesServices GamesServices
	AuthServices  AuthServices
}

func (uih *UserInteractionHandler) LikeGame(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)

	claims := user.Claims.(*services.JwtCustomClaims)

	idUser, err := uih.AuthServices.GetUserId(claims.Email)

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return renderView(c, errors_pages.Error500Index())
	}

	if c.Request().Method == "POST" {
		id := c.FormValue("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			return renderView(c, errors_pages.Error400Index())
		}

		err = uih.GamesServices.LikeGameByID(idInt, idUser)

		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return renderView(c, errors_pages.Error500Index())
		}

		return renderView(c, gamesviews.LikeButton(idInt, true))
	}

	games, err := uih.GamesServices.GetGamesLikedByUser(idUser)

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return renderView(c, errors_pages.Error500Index())
	}

	return renderView(c, gamesviews.GameIndexLiked(games))
}
