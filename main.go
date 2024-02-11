package main

import (
	"flag"
	"os"

	"github.com/Nuxnuxx/showcase_go/internal/database"
	"github.com/Nuxnuxx/showcase_go/internal/handlers"
	"github.com/Nuxnuxx/showcase_go/internal/services"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	PORT := flag.String("port", ":" + os.Getenv("PORT"), "port to run the server on")

	store, err := database.NewStore(os.Getenv("DB_NAME"))

	if err != nil {
		e.Logger.Fatal(err)
	}

	gameServices := services.NewGamesServices(services.Game{}, store, os.Getenv("API_KEY"))
	gameHandler := handlers.NewGamesHandlers(gameServices)

	handlers.SetupRoutes(e, gameHandler)

	// Start the server
	e.Logger.Fatal(e.Start(*PORT))
}
