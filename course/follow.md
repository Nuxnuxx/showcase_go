# Scaffolding

```bash
go mod init github.com/user/project
```

After that create a ```main.go``` and add it.
```go
// Filename: main.go
package main

func main(){

}
```

This is the entrypoint of your binary.

# Setting up

## Intall echo

```bash
go get github.com/labstack/echo/v4
```

## Get your Api Key for RAWG

What is RAWG : [Here](https://rawg.io/)

[Get your api key here](https://rawg.io/apidocs)

Now put it in your .env.
```Makefile
//Filename: .env
API_KEY="yourapikey"
```

## Load .env

This package his self explanatory, it just autoload what is in the ```.env``` automatically for us.

```bash
go get "github.com/joho/godotenv"
```

```go
// Filename: main.go
import (
    ...
    _ "github.com/joho/godotenv/autoload"
)
```

## Get started with RAWG

```go
// Filename: main.go

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
```

## Lets run this

```Makefile
// Filename: Makefile
build:
		@go build -o bin/app

run: build
		@./bin/app

// Nobody write test
test:
		@go test -v ./...
```

```make run``` it will build and then start the binary.

to see a more clear response you select only the name of the games with ```jq```.

```bash
make run | jq -r '.results[].name'
```

### | Checkpoint |
```bash
git reset --hard HEAD
```
```bash
git merge origin/get_started_with_rawg
```
### | Checkpoint |


# Getting started for real

## Layout Structure
```
├── main.go (entrypoint of our app)
├── go.sum // package.json
├── go.mod // package-lock.json
├── Makefile
├── .env
├── .gitignore
├── .git
├── internal // a special directory recognised which will prevent one package from being imported by another unless both share a common ancestor
│   ├── assets
│   │   ├── css
│   │   ├── js
│   ├── views // templates
│   ├── database // migration and factory for the database
│   ├── services // abstraction layer of database
│   ├── handlers // business logic
```

## Database

First we need to install sqlite driver for golang.
```bash
go get github.com/mattn/go-sqlite3
```

Then we can create the database file.
```go
//Filename: internal/database/database.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Db *sql.DB
}

func NewStore(dbName string) (Store, error) {
	Db, err := getConnection(dbName)
	if err != nil {
		return Store{}, err
	}

	if err := createMigrations(dbName, Db); err != nil {
		return Store{}, err
	}

	return Store{
		Db,
	}, nil
}

func getConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	if db != nil {
		return db, nil
	}

	// Init SQLite3 database
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		// log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}

func createMigrations(dbName string, db *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(64) NOT NULL
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
```

- Because this project is more focus on the big picture of the project we will not write everything our hand

### Factory pattern

In golang the most use pattern is the factory pattern which you already see without knowing here.
```go
type Store struct {
	Db *sql.DB
}

func NewStore(dbName string) (Store, error) {
	Db, err := getConnection(dbName)
	if err != nil {
		return Store{}, err
	}

	if err := createMigrations(dbName, Db); err != nil {
		return Store{}, err
	}

	return Store{
		Db,
	}, nil
}
```
#### Explanation

- Instead of calling the new keyword which you would do in typical language such as java, here we create a function that instantiate the object for us

```java
MyClass obj = new MyClass(10);
```

- And with factory
```java
Product product = ProductFactory.createProduct();
```

- In go, a factory is just a function or method that return an instance of a particular struct or interface

- Advantages :
    - Encapsulation: It encapsulates the object creation logic, so the client code doesn't need to know how objects are created.
    - Flexibility: It allows you to change the implementation of the object creation without affecting the client code.
    - Abstraction: It promotes loose coupling by allowing the client code to depend on abstractions rather than concrete implementations.

[If you want to know more about it](https://blog.matthiasbruns.com/golang-factory-method-pattern)

## Server

First of all, we still have our test about RAWG in your ```main.go```,you can get rid of everything in the main function and the useless import, let's get started with a real server.
```go
//Filename: main.go
func main() {
	e := echo.New()

	PORT := flag.String("port", ":" + os.Getenv("PORT"), "port to run the server on")

	// Start the server
	e.Logger.Fatal(e.Start(*PORT))

}
```

modify your .env to give it a PORT.
```Makefile
//Filename: .env
...
PORT="8080"
```

You can try with ```make run ``` it should start a echo server on the chosen port.

Now we can instantiate your store (database) in the entrypoint of our app.
```go
//Filename: main.go
PORT := ....

store, err := database.NewStore(os.Getenv("DB_NAME"))

if err != nil {
    e.Logger.Fatal(err)
}
```

And do the same as before in your ```.env```.

## Services

Now that we have a database we can create your first services which will be the games services, here the types needed for this services (thanks to gpt):
```go 
//Filename: internal/services/games.services.go
type GameService struct {
	Game Game
	GameStore database.Store
    ApiKey string
}

type EsrbRating struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Platform struct {
	ID           int    `json:"id"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	ReleasedAt   string `json:"released_at"`
	Requirements struct {
		Minimum     string `json:"minimum"`
		Recommended string `json:"recommended"`
	} `json:"requirements"`
}

type Game struct {
	ID               int         `json:"id"`
	Slug             string      `json:"slug"`
	Name             string      `json:"name"`
	Released         string      `json:"released"`
	Tba              bool        `json:"tba"`
	BackgroundImage  string      `json:"background_image"`
	Rating           int         `json:"rating"`
	RatingTop        int         `json:"rating_top"`
	Ratings          interface{} `json:"ratings"`
	RatingsCount     int         `json:"ratings_count"`
	ReviewsTextCount string      `json:"reviews_text_count"`
	Added            int         `json:"added"`
	AddedByStatus    interface{} `json:"added_by_status"`
	Metacritic       int         `json:"metacritic"`
	Playtime         int         `json:"playtime"`
	SuggestionsCount int         `json:"suggestions_count"`
	Updated          string   `json:"updated"`
	EsrbRating       EsrbRating  `json:"esrb_rating"`
	Platforms        []Platform  `json:"platforms"`
}

type Response struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Game `json:"results"`
}
```

After that we can create our function to initiate the service and follow the ```factory pattern``` for this service.
```go
//Filename: internal/services/games.services.go
func NewGamesServices(g Game, gStore database.Store, apiKey string) *GameService{

	return &GameService{
		Game:      g,
		GameStore: gStore,
		ApiKey: apiKey,
	}
}
```

For the first function we will create a function that get the games list.

```go
func (gs *GameService) GetGames(page int) ([]Game, error) {

	// Make the url
	builder := strings.Builder{}
	builder.WriteString("https://api.rawg.io/api/games?key=")
	builder.WriteString(os.Getenv("API_KEY"))

	// If page is not the first page, add the page number to the requestt 
	if (page > 0) {
		builder.WriteString(fmt.Sprintf("&page=%d", page))
	}

	// Make the request
	resp, err := http.Get(builder.String())

	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}

	defer resp.Body.Close()


	// This part bind the response to the struct
	var response Response

	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return response.Results, nil
}
```

The most important here is when we bind the reponse to the struct at the end, it is a common pattern in golang, after the data has been bind to the struct we have all the advantages of a strongly typed data.

### Test it

What's better to test your new services but to do a test (yes a test !) so we start by creating a new files in the folder ```services``` named ```game.services_test.go```.

After that we need a new package to assert some data.
```bash
go get github.com/stretchr/testify/assert
```

Let's write our test for the new function we just created.

```go
//Filename: internal/services/game.service_test.go
func TestGetGamesByPage(t *testing.T) {
	// Arrange
	store, err := database.NewStore("test.db")

	if err != nil {
		t.Fatalf("🔥 failed to connect to the database: %s", err)
	}

	t.Setenv("API_KEY", "yourapikey")

	gameService := NewGamesServices(Game{}, store, os.Getenv("API_KEY"))

	// Act
	result, err := gameService.GetGamesByPage(0)

	if err != nil {
		t.Fatalf("🔥 failed to get games: %s", err)
	}

	// Assert
	assert.NotNil(t, result)
}
```

## Handlers

In the handlers package we will use the services we already created to make some routes.

First of all create ```games.handlers.go``` in the handlers folder in it we wills store all the endpoint that in relation with games.


```go
//Filename: internal/handlers/games.handlers.go
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
```

Did you recognize some pattern ?

Yes first we have the factory pattern that we already encountered before and a new one.

### Dependency injection

First what is dependency injection ?

- It is the last letter of the famous acronym S.O.L.I.D, the D
- and it allow to switch the implementation of some dependency which means more simple unit testing because you can mock it, and even change it at the runtime ! 

It is exactly what happens here we create a interface (contract) and then we can pass any implementation that approve the contract.

```go
//Filename: internal/handlers/games.handlers.go

// The contract for the implementation
type GamesServices interface {
	GetGamesByPage(page int) ([]services.Game, error)
}

// Any implementation as paramaters that validate the contract
func NewGamesHandlers(gs GamesServices) *GamesHandler
```

Now that we have setting up your new handlers let's create your first real handler.
```go
//Filename: internal/handlers/games.handlers.go
func (gh *GamesHandler) GetGamesByPage(c echo.Context) error {

    // strconv use to transform string to int
    // because everything that come from web is string typed
	page, err := strconv.Atoi(c.QueryParam("page"))


    // if error when converting we return Bad request
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page")
	}

    // We use the service we created before
	games, err := gh.GamesServices.GetGamesByPage(page)


    // if error appears in services means it is a server error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(200, games)
}
```

The handler is standalone and unuse in our server we need to link it to the server, first create a ```routes.go``` in the handlers folder and add a new function.

```go
//Filename: internal/handlers/routes.go
func SetupRoutes(e *echo.Echo, gh *GamesHandler) {
	e.GET("/", gh.GetGamesByPage)
}
```

This function will be your router for the whole app.


We can use it in the entrypoint of the app to register this new handler

```go
//Filename: main.go
if err != nil {
    ...
}

gameServices := services.NewGamesServices(services.Game{}, store, os.Getenv("API_KEY"))
gameHandler := handlers.NewGamesHandlers(gameServices)

handlers.SetupRoutes(e, gameHandler)

//Start the server
```

Now if you can use the new endpoint we just created simply we can ```make run``` and go to this link

- [Here](http://localhost:8080?page=0) or copy this below
```
http://localhost:8080?page=0
```

## But the views folder, how can i be an HTML engineer without it

Yes i forgot HTML engineer here we are, first lets start by adding some new package

```bash
go install github.com/a-h/templ/cmd/templ@latest // for the cli
go get github.com/a-h/templ // for the code
```
the doc is [here](https://templ.guide/quick-start/installation/) if needed


It is a template engine which is staticly typed it will help keep the codebase clear

the extension ```.templ``` don't change the behavior of go we can use it as a normal go file like we would normally do

So let's jump to our views folder and create a new folder in it call ```layout``` in which we will create a new file call ```base.layout.templ``` and write some HTML

```go
//Filename: internal/views/layout/base.layout.templ
// dont miss it
package layout

templ Base() {
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<meta name="description" content="Lets go HTML Engineer"/>
			<meta name="google" content="notranslate"/>
			<link rel="stylesheet" href="/css/styles.css"/>
			<title>Games App</title>
            //Link for daisyui (dont use it in production)
			<link href="https://cdn.jsdelivr.net/npm/daisyui@4.6.2/dist/full.min.css" rel="stylesheet" type="text/css"/>
            // Link to tailwindUI (dont use it in production)
			<script src="https://cdn.tailwindcss.com"></script>
		</head>
		<body>
            { children... }
		</body>
	</html>
}
```

This will be the root of our website that we will integrate in it the navbar,main and also footer

The interesting part is the ```children...``` which means that we can use it as a layout and everything we pass in will be at this place, better than word. Let's see by example:

Now we need to have views for our games let get started with a view to show a list of games, create a new folder in ```views``` which we will call ```games_views``` and add a file ```game.list.templ```


```go
//Filename: internal/views/games_views/game.list.templ
package gamesviews 

templ GameCard(game services.Game){
	<div>
		<h2>{game.Name}</h2>
		<p>{game.Released}</p>
	</div>
}

templ GamesList(games []services.Game){
	// a for loop from golang
	for _, game := range games{
		@GameCard(game)
	}
}

templ GameIndex(games []services.Game){
	@layout.Base(){
		<h1>Games</h1>
		@GamesList(games)
	}
}
```

Here the interesting parts, first the usage of ```layout.Base``` in the ```GameIndex``` and the for loop which looks exactly the same as in go, if you know go you know templ.


We have finish for the folder ```views``` for the moment, we can now adapt our code to serve HTML and not JSON.

Before that, to generate the template you need to run ```templ generate```.

Now in the handlers for the game list we can change some lines and add a new files ```utils.go``` which will be all your utility function to render.

```go
//Filename: internal/handlers/utils.go

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

```

This function serve one purpose, elimination all the boilerplate to render a templ template, we don't need to know every aspect of it.
All we care is we pass a echo.Context and a component and it is render.

Now modify the handlers for the game list
```go
//Filename: internal/handlers/game.handlers.go

func (gh *GamesHandler) GetGamesByPage(c echo.Context) error {
    ...

    return c.JSON(200, games) -> return renderView(c, gamesviews.GameIndex(games))
}
```

```make run``` and check on the same page we go earlier.

You should see a ***beautiful*** list (joking) but it works, we have render our first HTML page.

### | Checkpoint |
```bash
git reset --hard HEAD
```
```bash
git merge origin/echo_server
```
### | Checkpoint |
