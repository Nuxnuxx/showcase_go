package handlers

import "github.com/labstack/echo/v4"


func SetupRoutes(e *echo.Echo, gh *GamesHandler) {
	e.GET("/", gh.GetGamesByPage)
}
