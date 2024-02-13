package handlers

import (
	"github.com/Nuxnuxx/showcase_go/internal/services"
	"github.com/Nuxnuxx/showcase_go/internal/views/layout"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, gh *GamesHandler, as *AuthHandler) {
	e.GET("/", HomeHandler)
	e.GET("/list", gh.GetGamesByPage)
	e.GET("/game/:id", gh.GetGameById)

	e.GET("/register", as.Register)
	e.POST("/register", as.Register)
	e.GET("/login", as.Login)

	protectedRoute := e.Group("/protected", echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JwtCustomClaims)
		},
		SigningKey:  as.AuthServices.GetSecretKey(),
		TokenLookup: "cookie:user",
	}))

	protectedRoute.GET("/profil", as.Profil)
}


func HomeHandler(c echo.Context) error {
	return renderView(c, layout.HomeIndex())
}
