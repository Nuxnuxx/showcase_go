package handlers

import (
	"github.com/Nuxnuxx/showcase_go/internal/views/layout"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, gh *GamesHandler, as *AuthHandler) {
	e.GET("/", HomeHandler)
	e.GET("/list", gh.GetGamesByPage)
	e.GET("/game/:id", gh.GetGameById)

	e.GET("/register", as.Register)
	e.POST("/register", as.Register)
}


func HomeHandler(c echo.Context) error {
	return renderView(c, layout.HomeIndex())
}
