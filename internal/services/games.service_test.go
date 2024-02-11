package services

import (
	"os"
	"testing"

	"github.com/Nuxnuxx/showcase_go/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestGetGamesByPage(t *testing.T) {
	// Arrange
	store, err := database.NewStore("test.db")

	if err != nil {
		t.Fatalf("🔥 failed to connect to the database: %s", err)
	}

	t.Setenv("API_KEY", "da7d0a35423b4557abcac0684875a989")

	gameService := NewGamesServices(Game{}, store, os.Getenv("API_KEY"))

	// Act
	result, err := gameService.GetGamesByPage(0)

	if err != nil {
		t.Fatalf("🔥 failed to get games: %s", err)
	}

	// Assert
	assert.NotNil(t, result)
}
