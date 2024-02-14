package main

import (
	"flag"
	"os"

	"github.com/Nuxnuxx/showcase_go/internal/database"
	"github.com/Nuxnuxx/showcase_go/internal/handlers"
	"github.com/Nuxnuxx/showcase_go/internal/services"
	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Static("/internal/css", "internal/css")

	PORT := flag.String("port", ":"+os.Getenv("PORT"), "port to run the server on")

	store, err := database.NewStore(os.Getenv("DB_NAME"))

	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Validator = &services.CustomValidator{Validator: validator.New()}

	gameServices := services.NewGamesServices(services.Game{}, store, os.Getenv("API_KEY"))
	gameHandler := handlers.NewGamesHandlers(gameServices)

	authServices := services.NewAuthServices(services.User{}, store, os.Getenv("SECRET_KEY"))
	authHandler := handlers.NewAuthHandler(authServices)

	userInteractionHandler := handlers.NewUserInteractionHandler(authServices, gameServices)

	handlers.SetupRoutes(e, gameHandler, authHandler, userInteractionHandler)

	// Start the server
	e.Logger.Fatal(e.Start(*PORT))
}
