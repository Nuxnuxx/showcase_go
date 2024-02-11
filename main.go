package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// GET API_KEY from .env file
	API_KEY := os.Getenv("API_KEY")

	// Build the URL
	builder := strings.Builder{}
	builder.WriteString("https://api.rawg.io/api/games?key=")
	builder.WriteString(API_KEY)

	// Make the request
	resp, err := http.Get(builder.String())

	if err != nil {
		panic(err)
	}

	// Defer the closing of the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	// Print the response body
	fmt.Println(string(body))

	// Now the body is close automatically because there is no need for it anymore
}
