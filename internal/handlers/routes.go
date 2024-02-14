package handlers

import (
	"github.com/Nuxnuxx/showcase_go/internal/services"
	"github.com/Nuxnuxx/showcase_go/internal/views/layout"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, gh *GamesHandler, ah *AuthHandler, uih *UserInteractionHandler) {
	e.GET("/", HomeHandler)
	e.GET("/list", gh.GetGamesByPage)
	e.GET("/game/:id", gh.GetGameById)

	authRouter := e.Group("/auth")
	authRouter.POST("/logout", ah.Logout)
	authRouter.Use(ah.CheckLogged)
	authRouter.GET("/register", ah.Register)
	authRouter.POST("/register", ah.Register)
	authRouter.GET("/login", ah.Login)
	authRouter.POST("/login", ah.Login)

	protectedRoute := e.Group("/protected", echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JwtCustomClaims)
		},
		ErrorHandler: ah.CheckNotLogged,
		SigningKey:   ah.AuthServices.GetSecretKey(),
		TokenLookup:  "cookie:user",
	}))

	protectedRoute.GET("/profil", ah.Profil)
	protectedRoute.GET("/liked", uih.LikeGame)
	protectedRoute.POST("/liked", uih.LikeGame)
}

func HomeHandler(c echo.Context) error {
	return renderView(c, layout.HomeIndex())
}
