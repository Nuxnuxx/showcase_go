# Scaffolding

```bash
go mod init github.com/user/project
```

After that create a `main.go` and add it.

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

This package his self explanatory, it just autoload what is in the `.env` automatically for us.

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

`make run` it will build and then start the binary.

to see a more clear response you select only the name of the games with `jq`.

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

First of all, we still have our test about RAWG in your `main.go`,you can get rid of everything in the main function and the useless import, let's get started with a real server.

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

You can try with `make run ` it should start a echo server on the chosen port.

Now we can instantiate your store (database) in the entrypoint of our app.

```go
//Filename: main.go
PORT := ....

store, err := database.NewStore(os.Getenv("DB_NAME"))

if err != nil {
    e.Logger.Fatal(err)
}
```

And do the same as before in your `.env`.

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

After that we can create our function to initiate the service and follow the `factory pattern` for this service.

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

What's better to test your new services but to do a test (yes a test !) so we start by creating a new files in the folder `services` named `game.services_test.go`.

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

First of all create `games.handlers.go` in the handlers folder in it we wills store all the endpoint that in relation with games.

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

The handler is standalone and unuse in our server we need to link it to the server, first create a `routes.go` in the handlers folder and add a new function.

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

Now if you can use the new endpoint we just created simply we can `make run` and go to this link

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

the extension `.templ` don't change the behavior of go we can use it as a normal go file like we would normally do

So let's jump to our views folder and create a new folder in it call `layout` in which we will create a new file call `base.layout.templ` and write some HTML

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

The interesting part is the `children...` which means that we can use it as a layout and everything we pass in will be at this place, better than word. Let's see by example:

Now we need to have views for our games let get started with a view to show a list of games, create a new folder in `views` which we will call `games_views` and add a file `game.list.templ`

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

Here the interesting parts, first the usage of `layout.Base` in the `GameIndex` and the for loop which looks exactly the same as in go, if you know go you know templ.

We have finish for the folder `views` for the moment, we can now adapt our code to serve HTML and not JSON.

Before that, to generate the template you need to run `templ generate`.

Now in the handlers for the game list we can change some lines and add a new files `utils.go` which will be all your utility function to render.

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

`make run` and check on the same page we go earlier.

You should see a **_beautiful_** list (joking) but it works, we have render our first HTML page.

### | Checkpoint |

```bash
git reset --hard HEAD
```

```bash
git merge origin/echo_server
```

### | Checkpoint |

# Now let's make it look like a real website

For a goodlooking website we need a navbar, footer those things are called `partials`.

So start by creating a folder call `partials` in the views folder in which we add a file call `navbar.partial.templ`.

**_EXERCICE_**

This navbar should be able to go to page `/`,`/list`,`profil` and `/liked` which are self explanatory.

Enjoy and try do something cool, then implement it to the layout

**_CORRECTION_**

```go
//Filename: internal/views/partials/navbar.partial.templ

templ NavBar(){
	<nav class="navbar bg-primary text-primary-content fixed top-0 z-10">
		<div class="navbar-start">
			<a hx-swap="transition:true" class="btn btn-ghost text-xl" href="/">
				Todo List
			</a>
		</div>
		<div class="navbar-end">
				<a class="btn btn-ghost text-lg" href="/">
					Home
				</a>
				<a class="btn btn-ghost text-lg" href="/list">
					List
				</a>
				<a class="btn btn-ghost text-lg" href="/liked">
					Liked
				</a>
				<a class="btn btn-ghost text-lg" href="/profil">
					Profil
				</a>
		</div>
	</nav>
}

//Filename: internal/views/layout/base.layout.templ
<body>
    <header>
        @partials.NavBar()
    </header>
    { children... }
</body>
```

Ok now that you have played a little with the templating engine, don't you find annoying to need to rerun `templ generate` everytime you make a change, let's fix this go in your `Makefile`

```Makefile
build:
    @templ generate
```

Now it will recreate template at every rerun of `make run`

You can also make a footer if you want !

For now all of this is just frontend stuff and no ones of those link are working, let's get on it

First the home page we need a page that informs the user of what is the website

Create a new file in the `layout` folder call `homepage.layout.templ` and insert your homepage templ

```go
//Filename: internal/views/layout/homepage.layout.templ
templ Home() {
	<div class="container mx-auto mt-8">
		<section class="text-center">
			<h1 class="text-4xl font-bold text-gray-800 mb-4">Welcome to GameApp!</h1>
			<p class="text-lg text-gray-600 mb-8">Discover and like your favorite games.</p>
			<a href="/login" class="bg-purple-500 text-white px-6 py-3 rounded-md text-lg hover:bg-blue-600">Get Started</a>
		</section>
	</div>
}

templ HomeIndex() {
	@Base() {
		@Home()
	}
}
```

Then add a new endpoint in `routes.go` and modify the current one

```go
//Filename: internal/handlers/routes.go
func SetupRoutes(e *echo.Echo, gh *GamesHandler) {
	e.GET("/", HomeHandler)
	e.GET("/list", gh.GetGamesByPage)
}

func HomeHandler(c echo.Context) error {
	return renderView(c, layout.HomeIndex())
}
```

Here we have modify the endpoint to get the game of list on the path `/list` and create a new handler in which we render our new views

The list endpoint don't work anymore because if we don't pass a query parameters `page=x` it crash and say invalid page why ?

Because of this part

```go
//Filename: internal/handlers/games.handlers.go
func (gh *GamesHandler) GetGamesByPage(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))

    // This part just return JSON
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid page")
	}

	games, err := gh.GamesServices.GetGamesByPage(page)

    // This is also handle with a return as a JSON
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong")
	}


	return renderView(c, gamesviews.GameIndex(games))
}
```

So we need a way to show error in a more pretty / web app flavor, one of the common way to do it is to create a view for each error or redirect to another page.

We could also just change the behavior of this handler to just send the page 0 if it had receive no page on his parameters.

Finally, we will do both, first change the function to not crash when there is no query params.

```go
page := c.QueryParam("page")

if page == "" {
    page = "0"
}

pageInt, err := strconv.Atoi(page)

if err != nil {
    return c.JSON(http.StatusBadRequest, "Invalid page")
}
```

Now the start of the handler should look like this, but we are still handling error with JSON response, let's fix it.

Create a new folders in `views` named `errors_pages` and add one named `error.400.templ` in it we will show a more looking good errors

```go
//Filename: internal/views/errors_page/error.400.templ
templ Error400() {
	<section class="flex flex-col items-center justify-center h-[100vh] gap-4">
		<div class="items-center justify-center flex flex-col gap-4">
			<h1 class="text-9xl font-extrabold text-gray-700 tracking-widest">
				400
			</h1>
			<h2 class="bg-rose-700 px-2 text-sm rounded rotate-[20deg] absolute">
				Bad Request
			</h2>
		</div>
		<p class="text-xs text-center md:text-sm text-gray-400">
			Your request was malformed
		</p>
	</section>
}


templ Error400Index(){
	@layout.Base(){
		@Error400()
	}
}
```

This should looks good, now we can use it in the handler

```go
//Filename: internal/handlers/games.handlers.go

func (gh *GamesHandler) GetGamesByPage(c echo.Context) error {
	page := c.QueryParam("page")

	if page == "" {
		page = "0"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return renderView(c, errors_pages.Error400Index())
	}

	games, err := gh.GamesServices.GetGamesByPage(pageInt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Something went wrong")
	}


	return renderView(c, gamesviews.GameIndex(games))
}
```

**_EXERCICE_**
You can do the same for 500 now.

Next we should redesign this horrible list page.

```go
//Filename: internal/views/games_views/game.list.templ
templ GameCard(game services.Game){
    <div class="card">
			<figure><img class="object-cover max-h-60 max-w-full" src={game.BackgroundImage} alt={game.Name} /></figure>
			<div class="card-body">
					<h2 class="card-title">{game.Name}</h2>
					<p>{game.Released}</p>
			</div>
    </div>
}

templ GamesList(games []services.Game){
    <div class="grid grid-cols-3 gap-4">
        // Loop through games
        for _, game := range games{
            <div class="col">
                @GameCard(game)
            </div>
        }
    </div>
}

templ GameIndex(games []services.Game){
    @layout.Base(){
        <h1>Games</h1>
        @GamesList(games)
    }
}

```

### Here HTMX Come

Now that it looks fine, let's focus on a really cool features `infinite scroll`, here HTMX come to play

[Here](https://htmx.org/examples/infinite-scroll/) is an example of a infinite scroll in HTMX, look pretty simple isn't it.

First let's install htmx, add this to the `header` of the layout template.

```go
//Filename: internal/views/layout/base.layout.templ
<script src="https://unpkg.com/htmx.org@1.9.10" integrity="sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC" crossorigin="anonymous"></script>
```

First we need to pass the currentPage from the backend to the frontend by changing the handler

```go
//Filename: internal/handlers/games.handlers.go
return renderView(c, gamesviews.GameIndex(games)) -> return renderView(c, gamesviews.GameIndex(games, pageInt))
```

Then modifying

```go
//Filename: internal/views/games_views/game.list.templ

templ GamesList(games []services.Game, currentPage int) {
	<div id="game_list" class="grid grid-cols-3 gap-4">
		// Loop through games
		for i, game := range games {
			if i == len(games) - 1 {
				<div
					id="load_more"
					hx-trigger="revealed"
					hx-get={ "/list?page=" + strconv.Itoa(currentPage+1) }
				>
					@GameCard(game)
				</div>
			}
			@GameCard(game)
		}
	</div>
}

templ GameIndex(games []services.Game, currentPage int) {
	@layout.Base() {
		@GamesList(games, currentPage)
	}
}
```

But what happens ! Our result has been in one card only but why ?

Because by default HTMX replace the element from which you make the request we need to specify which [target](https://htmx.org/attributes/hx-target/) and how we want to [swap](https://htmx.org/attributes/hx-swap/) them.

Also it will bug because it also send the `<div id="game_list" class="grid grid-cols-3 gap-4">` around it for each request we need to move it to the index

We also need to send a new GameList when the page is more then 0.

Which give us this

```go
//Filename: internal/views/games_views/game.list.templ

templ GamesList(games []services.Game, currentPage int) {
	// Loop through games
	for i, game := range games {
		if i == len(games) - 1 {
			<div id="load_more" hx-trigger="revealed" hx-get={ "/list?page=" + strconv.Itoa(currentPage+1) } hx-target="#game_list" hx-swap="beforeend">
				@GameCard(game)
			</div>
		}
		@GameCard(game)
	}
}

templ GameIndex(games []services.Game, currentPage int) {
	@layout.Base() {
		<div id="game_list" class="grid grid-cols-4 gap-4">
			@GamesList(games, currentPage)
		</div>
	}
}
```

Add it at the end of the end of the handlers before the first render return.

```go
//Filename: internal/handlers/games.handlers.go

if pageInt > 0 {
        return renderView(c, gamesviews.GamesList(games, pageInt))
}

```

Also we can add `hx-boost="once"` so it only happens once and we can scroll back without causing a new request.

And boom we got infinite scroll with few lines.

Now that we are here we can also put `hx-boost="true` to the body so the anchor tag are now doing ajax request and not a full reload of the web app.

### Some find tuning for you developement experience

It would be cool not to have the need to run make run at every changes, that when [air](https://github.com/cosmtrek/air) comes along, with this binary we can define a config file in which it can run every command we want after certain files has been changed.

First install air by going on their [github](https://github.com/cosmtrek/air).

Run `air init`, you should have now a `.air.toml` created.

Let's modify the config files, first add a new extension files to watch one

```toml
include_ext = ["go", "tpl", "tmpl", "html", "templ"]
```

And modify the bin and command to run on every reload

```toml
bin = "./bin/app"
cmd = "make build"
```

Now we are fine if you run `air`, it should restart at every changes.

### The game detail page

First we need to make a new function in our services for the game, let's create a function to retrieve a game by id.

```go
//Filename: internal/services/games.services.go

func (gs *GameService) GetGamesByID(id int) (GameFullDetail, error){
	// Make the url
	builder := strings.Builder{}
	builder.WriteString("https://api.rawg.io/api/games/")
	builder.WriteString(strconv.Itoa(id))
	builder.WriteString("?key=")
	builder.WriteString(gs.ApiKey)

	resp, err := http.Get(builder.String())

	if err != nil {
		return GameFullDetail{}, fmt.Errorf("Error making request: %v", err)
	}

	defer resp.Body.Close()

	// This part bind the response to the struct
	var response GameFullDetail

	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &response); err != nil {
		return GameFullDetail{}, fmt.Errorf("Error unmarshalling response: %v", err)
	}

	return response, nil
}
```

We can't return a nil this time because the type `Game` is not a pointer such as below ere it was a array.

And a new type which we will be bind to it

```go
//Filename: internal/services/games.services.go

type GameFullDetail struct {
    ID                    int               `json:"id"`
    Slug                  string            `json:"slug"`
    Name                  string            `json:"name"`
    NameOriginal          string            `json:"name_original"`
    Description           string            `json:"description"`
    Metacritic            int               `json:"metacritic"`
    MetacriticPlatforms   []MetacriticPlatform `json:"metacritic_platforms"`
    Released              string            `json:"released"`
    TBA                   bool              `json:"tba"`
    Updated               string            `json:"updated"`
    BackgroundImage       string            `json:"background_image"`
    BackgroundImageAdditional string         `json:"background_image_additional"`
    Website               string            `json:"website"`
    Rating                float32              `json:"rating"`
    RatingTop             int               `json:"rating_top"`
    Reactions             map[string]interface{} `json:"reactions"`
    Added                 int               `json:"added"`
    AddedByStatus         map[string]interface{} `json:"added_by_status"`
    Playtime              int               `json:"playtime"`
    ScreenshotsCount      int               `json:"screenshots_count"`
    MoviesCount           int               `json:"movies_count"`
    CreatorsCount         int               `json:"creators_count"`
    AchievementsCount     int               `json:"achievements_count"`
    ParentAchievementsCount int          `json:"parent_achievements_count"`
    RedditURL             string            `json:"reddit_url"`
    RedditName            string            `json:"reddit_name"`
    RedditDescription     string            `json:"reddit_description"`
    RedditLogo            string            `json:"reddit_logo"`
    RedditCount           int               `json:"reddit_count"`
    TwitchCount           int            `json:"twitch_count"`
    YoutubeCount          int            `json:"youtube_count"`
    ReviewsTextCount      int            `json:"reviews_text_count"`
    RatingsCount          int               `json:"ratings_count"`
    SuggestionsCount      int               `json:"suggestions_count"`
    AlternativeNames      []string          `json:"alternative_names"`
    MetacriticURL         string            `json:"metacritic_url"`
    ParentsCount          int               `json:"parents_count"`
    AdditionsCount        int               `json:"additions_count"`
    GameSeriesCount       int               `json:"game_series_count"`
    ESRBRating            EsrbRating        `json:"esrb_rating"`
    Platforms             []Platform    `json:"platforms"`
}
```

After that we can create a handler to serve this page.

```go
//Filename: internal/handlers/games.handlers.go

func (gh *GamesHandler) GetGameById(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return renderView(c, errors_pages.Error400Index())
	}

	game, err := gh.GamesServices.GetGamesByID(idInt)

	if err != nil {
		fmt.Println(err)
		return renderView(c, errors_pages.Error500Index())
	}

	return renderView(c, gamesviews.GamePageIndex(game))
}
```

Don't forget to add our new function to the interface at the top of the file.

Add a view to show the information we just retrieve

```go
//Filename: internal/views/games_views/game.list.templ

templ GamePage(game services.GameFullDetail){
 <div class="max-w-4xl mx-auto shadow-md rounded-md p-6">
        <h1 class="text-3xl font-bold mb-4">{game.Name}</h1>
        <img src={game.BackgroundImage} alt={game.Name} class="w-full h-auto rounded-md mb-4"/>
        <h2 class="text-2xl font-bold mb-2">{game.Name}</h2>
        <p class="text-lg mb-4">Released: {game.Released}</p>
        <p class="text-base mb-4">{game.Description}</p>
				if (game.Website != ""){
					<a target="_blank" href={templ.URL(game.Website)} class="btn">{game.Website}</a>
				}
    </div>
}

templ GamePageIndex(game services.GameFullDetail){
	@layout.Base(){
		@GamePage(game)
	}
}
```

And what we need to do now ?

Yes add it to the routes

```go
//Filename: internal/handlers/routes.go

e.GET("/game/:id", gh.GetGameById)
```

# Authentification

For the authentification we will use [JsonWebToken](https://jwt.io/),

Jwt will be use here to store a json object on the client side which will be signed by a secret key and we will check that it is a valid token on back to do [Authentification](https://www.onelogin.com/learn/authentication-vs-authorization#:~:text=Authentication%20vs.-,Authorization,authorization%20determines%20their%20access%20rights.)

Echo provide a built-in middleware to get started with `JWT`.

```bash
go get github.com/labstack/echo-jwt/v4
```

After that we have all we need to our authentification, so we can create a new services called `auth.services.go` and start up it

```go
//Filename: internal/services/auth.services.go

func NewAuthServices(u User, uStore database.Store, secretKey string) *AuthService {

	return &AuthService{
		User:      u,
		UserStore: uStore,
		SecretKey: []byte(secretKey),
	}
}

type AuthService struct {
	User      User
	UserStore database.Store
	SecretKey []byte
}

type User struct {
	ID 		 int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
```

After that we can add the first important function which have to goal to create a new user

```go
//Filename: internal/services/auth.services.go

func (as *AuthService) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8) // Hash the password because we are not criminal
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(email, password, username) VALUES($1, $2, $3)`

	_, err = as.UserStore.Db.Exec(
		stmt,
		u.Email,
		string(hashedPassword),
		u.Username,
	)

	return err
}
```

Then we can create the handlers that come with it `auth.handlers.go` and start up it too

```go

type AuthServices interface {
	CreateUser(user services.User) error
}

func NewAuthHandler(as AuthServices) *AuthHandler {

	return &AuthHandler{
		AuthServices: as,
	}
}

type AuthHandler struct {
	AuthServices AuthServices
}
```

And add a handler that first show the form to the user

```go
//Filename: internal/services/auth.services.go

func (au *AuthHandler) Register(c echo.Context) error {
	return renderView(c, authviews.RegisterIndex())
}
```

Then create the view we use in the return, create a new folder in views called `auth_views` and a file named `auth.register.templ` then create the components has we always has done

```go
//Filename: internal/views/auth_views/auth.register.templ

templ Register() {
	<div class="bg-white p-8 rounded shadow-md w-full max-w-md">
		<h2 class="text-2xl font-semibold mb-4">User Registration</h2>
		<form action="" method="post">
			<div class="mb-4">
				<label for="email" class="block text-gray-700">Email:</label>
				<input type="email" id="email" name="email" required class="form-input mt-1 block w-full" />
			</div>
			<div class="mb-4">
				<label for="username" class="block text-gray-700">Username:</label>
				<input type="text" id="username" name="username" required class="form-input mt-1 block w-full" />
			</div>
			<div class="mb-4">
				<label for="password" class="block text-gray-700">Password:</label>
				<input type="password" id="password" name="password" required class="form-input mt-1 block w-full" />
			</div>
			<div class="mb-4">
				<button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Register</button>
			</div>
		</form>
	</div>
}

templ RegisterIndex() {
	@layout.Base() {
		@Register()
	}
}
```

From that we can add those new features to the routes by adding the new handlers and create new endpoint

```go
//Filename: internal/handlers/routes.go

func SetupRoutes(e *echo.Echo, gh *GamesHandler, as *AuthHandler) {
	e.GET("/", HomeHandler)
	e.GET("/list", gh.GetGamesByPage)
	e.GET("/game/:id", gh.GetGameById)

	e.GET("/register", as.Register)
}
```

Now it show an error in the main, that because we need to initiate the new service and handler to pass it to the setupRoutes function

```go
//Filename: main.go
gameHandler := ...

authServices := services.NewAuthServices(services.User{}, store, os.Getenv("SECRET_KEY"))
authHandler := handlers.NewAuthHandler(authServices)

handlers.SetupRoutes(e, gameHandler, authHandler)
```

And add a new SECRET_KEY in the `.env` file too.

For now our form doesn't do anything so let's make it functional, add a post endpoint to the setupRoutes

```go
//Filename: internal/handlers/routes.go

e.POST("/register", as.Register)
```

But what we use the same handler for the post and get, and yes we will handle the post send in the same, let's check how it works

```go
//Filename: internal/services/auth.services.go

func (au *AuthHandler) Register(c echo.Context) error {
	if c.Request().Method == "POST" {
		user := services.User{
			Email:    c.FormValue("email"),
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		err := au.AuthServices.CreateUser(user)

		if err != nil {
			return renderView(c, errors_pages.Error500Index())
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}
	return renderView(c, authviews.RegisterIndex())
}
```

Here we get the data we get from the form and pass it to a struct user to then create it with our function.

We still have some work to do on the view, just to add the endpoint in which we want to send the post request here `<form action="/register" method="post">`.

Now it works you should be redirect to the homepage.

But we still have to handle errors, use JWT to create protected routes and store the token on the client side.

### Handle errors on form

We will use the built-in [validator](https://echo.labstack.com/docs/request#validate-data) from echo to handle those, we need to add our constraints to the struct

Before that just need to install the validator framework

```bash
go get github.com/go-playground/validator
```

```go
//Filename: internal/services/auth.services.go

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}
```

This is self-explanatory, we check that the email is a valid email, and limit the size of username and password.

Now we create a custom validator in a new file `utils.go` in the services folder

```go
//Filename: internal/services/utils.go

type (
	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
```

Now we need to make the echo server aware of the new validator we just create.

```go
//Filename: main.go
e.Validator = &services.CustomValidator{Validator: validator.New()}

gameServices := ....
```

and now we can return a 400 if the value don't fill the constraint we set before

```go
//Filename: internal/services/auth.services.go
user := services.User{
    ////
}

if err := c.Validate(user); err != nil {
    return renderView(c, errors_pages.Error400Index()) }
```

But it is not enough and the user cannot know which error has cause it.

So now we have a error which we can cast to validateErrors object and use it to build the errors message to the user we need one function in the `utils.go` services file

```go
//Filename: internal/services/utils.go

type HumanErrors struct {
    Value 	 string
    Error 	 string
}

func CreateHumanErrors(err error) map[string]HumanErrors {
	errors := make(map[string]HumanErrors)

	for _, v := range err.(validator.ValidationErrors) {
		error := strings.Builder{}
		error.WriteString(fmt.Sprintf("%s should be %s %s",
			strings.Split(v.Namespace(), ".")[1],
			v.Param(),
			v.Tag(),
		))

		errors[strings.ToLower(v.Field())] = HumanErrors{
			Value: v.Value().(string),
			Error: error.String(),
		}
	}

	return errors
}
```

This function will help us get a map and not a simple error which give us more human readable error for the user and the old value if we need to put it back to the input.

And modify the inner if of validation to return the register components with the map of errors

```go
//Filename: internal/services/auth.services.go

if err := c.Validate(user); err != nil {
    humanErrors := services.CreateHumanErrors(err)

    return renderView(c, authviews.Register(humanErrors))
}
```

now we just need to modify some part of the view and it should work fine

```go
//Filename: internal/views/auth_views/auth.register.templ

templ Register(humanErrors map[string]services.HumanErrors) {
		<form hx-post="/register" hx-boost="true" hx-swap="outerHTML">
			<div class="mb-4">
				<label for="email" class="block text-gray-700">Email:</label>
				<input type="email" id="email" name="email" required class="form-input mt-1 block w-full" />
				if human, ok := humanErrors["email"]; ok {
					<div class="text-red-500 text-sm">{human.Error}</div>
				}
			</div>
			<div class="mb-4">
				<label for="username" class="block text-gray-700">Username:</label>
				<input type="text" id="username" name="username" required class="form-input mt-1 block w-full" />
				<div class="text-sm"> 3 - 20 characters, letters, numbers, and underscores only</div>
				if human, ok := humanErrors["username"]; ok {
					<div class="text-red-500 text-sm">{human.Error}</div>
				}
			</div>
			<div class="mb-4">
				<label for="password" class="block text-gray-700">Password:</label>
				<input type="password" id="password" name="password" required class="form-input mt-1 block w-full" />
				<div class="text-sm"> 8 - 50 characters, at least one letter, one number, and one special character</div>
				if human, ok := humanErrors["password"]; ok {
					<div class="text-red-500 text-sm">{human.Error}</div>
				}
			</div>
			<div class="mb-4">
				<button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Register</button>
			</div>
		</form>
}

templ RegisterIndex() {
	@layout.Base() {
		<div class="bg-white p-8 rounded shadow-md w-full max-w-md">
			<h2 class="text-2xl font-semibold mb-4">User Registration</h2>
			@Register(nil)
		</div>
	}
}
```

We use htmx to rerender only the form if there is an error and we can check validation as same as in other languages.

Now we can do all the part of `JWT` when we register in our app.

First we need to create the token and store it client-side after we are sure the user has really benn created.

We have forgot something, we don't check if the email already exist and it should cause a error 500 when creating a user that already exist, let's get on it.

```go
//Filename: internal/services/auth.services.go

func (as *AuthService) CheckEmail(email string) (User, error) {

	query := `SELECT email, password, username FROM users
		WHERE email = ?`

	stmt, err := as.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	as.User.Email = email
	err = stmt.QueryRow(
		as.User.Email,
	).Scan(
		&as.User.Email,
		&as.User.Password,
		&as.User.Username,
	)

	if err != nil {
		return User{}, err
	}

	return as.User, nil
}
```

And we now use it to check that the user already exist or not in our handlers.

```go

if err := c.Validate(user); err != nil {
    ////
}

userInDatabase, err := au.AuthServices.CheckEmail(user.Email)

if userInDatabase != (services.User{}) {
    humanErrors := map[string]services.HumanErrors{
        "email": {
            Error: "Email already exists",
            Value: user.Email,
        },
    }

    return renderView(c, authviews.Register(humanErrors))
}
```

Now the user will receive a error to show him that he already as an account.

```go
//Filename: internal/handlers/auth.handlers.go
func (au *AuthHandler) Register(c echo.Context) error {
    ////

    token, err := au.AuthServices.GenerateToken(user)

    if err != nil {
        return renderView(c, errors_pages.Error500Index())
    }

    cookie := http.Cookie{
        Name: "token",
        Value: token,
        Path:    "/",
        HttpOnly: true,
        Secure: true,
        Expires: time.Now().Add(24 * time.Hour),
    }

    c.SetCookie(&cookie)

    // INFO: To redirect when using HTMX you need to set the HX-Redirect header
    c.Response().Header().Set("HX-Redirect", "/")
    c.Response().WriteHeader(http.StatusOK)
    return nil
}
```

We also need the GenerateToken function we just use to create the token

```go
//Filename: internal/services/auth.services.go

type JwtCustomClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (as *AuthService) GenerateToken(user User) (string, error) {
	claims := &JwtCustomClaims{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	signedToken, err := token.SignedString(as.SecretKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
```

The token JWT work has defined here, we create claims which are the data that will contains the tokenthen we register a timestamp of now and and an expire time 24 hours later, to make sure that if the token leaks it will expire at a times, after that we sign it with our secret key.

Now that we have the cookies store client side we need to create a protected routes and check that the token is valid.

We will also need the secret key from the services of auth.

```go
//Filename: internal/services/auth.services.go

func (as *AuthService) GetSecretKey() string {
	return as.SecretKey
}
```

Also put in the interface of the handlers

```go
//Filename: internal/handlers/auth.handlers.go

type AuthServices interface {
	GetSecretKey() string
    ///
}
```

Now we can create our new routes which will be protected

```go
//Filename: internal/handlers/routes.go

protectedRoute := e.Group("/protected", echojwt.WithConfig(echojwt.Config{
    NewClaimsFunc: func(c echo.Context) jwt.Claims {
        return new(services.JwtCustomClaims)
    },
    SigningKey: as.AuthServices.GetSecretKey(),
    TokenLookup: "cookie:token",
}))

protectedRoute.GET("/", HomeHandler)
```

Here we pass the new claims type we just create, the signing key from the auth services to decrypt it, and where to find the token in the incoming request.

And you can also refacto all the place where you have a `return renderView(c, ErrorXXX)` to add just before the good status on the Response

```go
// Everywhere
c.Response().WriteHeader(correct status)
return renderView(/////)
```

Now we can do the most simple routes, a profil page. first modify the routes protected by that

```go
//Filename: internal/handlers/routes.go
protectedRoute.GET("/", HomeHandler) -> protectedRoute.GET("/profil", as.Profil)
```

Then we can create our new handler

```go
//Filename: internal/handlers/auth.handlers.go

func (au *AuthHandler) Profil (c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)

	if !ok {
		log.Errorf("Error getting claims from token: %v", token)
		return renderView(c, errors_pages.Error401Index())
	}

	claims, ok := token.Claims.(*services.JwtCustomClaims)

	if !ok {
		log.Errorf("Error getting claims from token: %v", token)
		return renderView(c, errors_pages.Error401Index())
	}

	user := services.User{
		Email:    claims.Email,
		Username: claims.Username,
	}

	return renderView(c, authviews.ProfilIndex(user))
}
```

Now we need to add a 401 page to the error pages.

After that we can also create our profilIndex.

```go
//Filename: internal/views/auth_views/auth.register.templ

templ Profil(user services.User){
 <div class="max-w-md mx-auto bg-white rounded-lg shadow-lg overflow-hidden">
    <div class="p-6">
      <h2 class="text-2xl font-semibold text-gray-800 mb-2">Profile Information</h2>
      <div class="space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
          <p class="text-lg font-semibold text-gray-900" id="username">{user.Username}</p>
        </div>
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
          <p class="text-lg font-semibold text-gray-900" id="email">{user.Email}</p>
        </div>
      </div>
    </div>
  </div>
}

templ ProfilIndex(user services.User){
	@layout.Base(){
		@Profil(user)
	}
}
```

And modify the navbar to point to the good endpoint and it will work fine.

So let's do the login page now, so create a new handler.

```go
//Filename: internal/handlers/auth.handlers.go

func (au *AuthHandler) Login(c echo.Context) error {
	return renderView(c, authviews.LoginIndex())
}
```

And create the form that come with it in a file `auth.login.templ` in the auth views folder.

```go
//Filename: internal/views/auth_views/auth.login.templ

templ Login(humanErrors map[string]services.HumanErrors) {
		<form hx-post="/register" hx-boost="true" hx-swap="outerHTML">
			<div class="mb-4">
				<label for="email" class="block text--700">Email:</label>
				<input type="email" id="email" name="email" required class="form-input mt-1 block w-full" />
				if human, ok := humanErrors["email"]; ok {
					<div class="text-red-500 text-sm">{human.Error}</div>
				}
			</div>
			<div class="mb-4">
				<label for="password" class="block text-white-700">Password:</label>
				<input type="password" id="password" name="password" required class="form-input mt-1 block w-full" />
				<div class="text-sm"> 8 - 50 characters, at least one letter, one number, and one special character</div>
				if human, ok := humanErrors["password"]; ok {
					<div class="text-red-500 text-sm">{human.Error}</div>
				}
			</div>
			<div class="mb-4">
				<button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Register</button>
			</div>
		</form>
}

templ LoginIndex() {
	@layout.Base() {
		<div class="p-8 rounded shadow-md w-full max-w-md mx-auto">
			<h2 class="text-2xl font-semibold mb-4">User Registration</h2>
			@Login(nil)
		</div>
	}
}
```

Then we declare it in the `setupRoutes` function.

```go
//Filename: internal/handlers/routes.go

e.GET("/login", as.Login)
```

Notice, you can go to login and register page, even if you are connected, let's fix that with a middleware.

```go
//Filename: internal/handlers/auth.handlers.go

func (au *AuthHandler) CheckAlreadyLogged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("user")

		if err != nil {
			fmt.Println(err)
			return next(c)
		}

		if token.Value != "" {
			fmt.Println(token)
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}
```

Now we need to include every routes that is use to connect with this middleware in the `setupRoutes`

```go
//Filename: internal/handlers/routes.go

authRouter := e.Group("/auth", as.CheckAlreadyLogged)
authRouter.GET("/register", as.Register)
authRouter.POST("/register", as.Register)
authRouter.GET("/login", as.Login)
```

And everywhere there is `/register` to redirect now to `/auth/register`.

We also need to redirect if the user is not connected and try to go profil or other routes that need an account to access it.

So we need to create a errror handler for the `echojwt.Config`

```go
//Filename: internal/handlers/auth.handlers.go

func (au *AuthHandler) CheckNotLogged(c echo.Context, err error) error {
	token, ok := c.Get("user").(string)

	// Means the user is not connected
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/auth/register")
	}

	// Means the user is not connected
	if token == "" {
		return c.Redirect(http.StatusSeeOther, "/auth/register")
	}

	return nil
}
```

And add it to the config.

```go
//Filename: internal/handlers/routes.go

protectedRoute := e.Group("/protected", echojwt.WithConfig(echojwt.Config{
    NewClaimsFunc: func(c echo.Context) jwt.Claims {
        return new(services.JwtCustomClaims)
    },
    ErrorHandler: as.CheckNotLogged,
    ////
}))
```

Now if we try to click on profil we will be redirect to the register page.

So le'ts make this login functional now.

First we can make a button `Already signed up` on the register page and we will do the form logic later on.

```html
//Filename: internal/views/auth_views/auth.register.templ

<div class="mb-4">
  <a
    href="/auth/login"
    class="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600"
    >Already registered</a
  >
</div>
```

Add it to the login Form, and now we can do the logic to connect.

```go
//Filename: internal/handlers/auth.handlers.go

if c.Request().Method == "POST" {
    user, err := ah.AuthServices.CheckEmail(c.FormValue("email"))

    // If the user is not found
    if err != nil {
        error := map[string]services.HumanErrors{
            "email": {
                Error: "Email not found",
                Value: c.FormValue("email"),
            },
        }

        return renderView(c, authviews.Login(error))
    }

    err = bcrypt.CompareHashAndPassword(
        []byte(user.Password),
        []byte(c.FormValue("password")),
    )

    // If the password is not correct
    if err != nil {
        error := map[string]services.HumanErrors{
            "internal": {
                Error: "Invalid Credentials",
                Value: "",
            },
        }

        return renderView(c, authviews.Login(error))
    }

    token, err := ah.AuthServices.GenerateToken(user)

    if err != nil {
        c.Response().WriteHeader(http.StatusInternalServerError)
        return renderView(c, errors_pages.Error500Index())
    }

    cookie := http.Cookie{
        Name:    "user",
        Value:   token,
        Path:    "/",
        Secure:  true,
        HttpOnly: true,
        Expires: time.Now().Add(24 * time.Hour),
    }

    c.SetCookie(&cookie)

    c.Response().Header().Set("HX-Redirect", "/")
    c.Response().WriteHeader(http.StatusOK)
    return nil
}
```

And add a POST request to the routes for the same handle as we have done with the register.

Now we modify the login template to follow what we have done.

```html
//Filename: internal/views/auth_views/auth.login.templ

<form hx-post="/auth/login" hx-boost="true" hx-swap="outerHTML">
  <div class="mb-4">
    <label for="email" class="block text--700">Email:</label>
    <input
      type="email"
      id="email"
      name="email"
      required
      class="form-input mt-1 block w-full"
    />
    if human, ok := humanErrors["email"]; ok {
    <div class="text-red-500 text-sm">{human.Error}</div>
    }
  </div>
  <div class="mb-4">
    <label for="password" class="block text-white-700">Password:</label>
    <input
      type="password"
      id="password"
      name="password"
      required
      class="form-input mt-1 block w-full"
    />
    <div class="text-sm">
      8 - 50 characters, at least one letter, one number, and one special
      character
    </div>
  </div>
  <div class="mb-4">
    <button
      type="submit"
      class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
    >
      Login
    </button>
  </div>
  if human, ok := humanErrors["internal"]; ok {
  <div class="text-red-500 text-sm">{human.Error}</div>
  }
</form>
```

we can also put a logout button on the profil page.

```html
//Filename: internal/views/auth_views/auth.register.templ

<div class="p-6">
  <form action="/auth/logout" method="post">
    <button
      type="submit"
      class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
    >
      Logout
    </button>
  </form>
</div>
```

And we need to add a endpoint to delete the cookie.

```go
//Filename: internal/handlers/routes.go

authRouter := e.Group("/auth")
authRouter.POST("/logout", as.Logout)
authRouter.Use(as.CheckLogged)
authRouter.GET("/register", as.Register)
authRouter.POST("/register", as.Register)
authRouter.GET("/login", as.Login)
authRouter.POST("/login", as.Login)
```

We dont check if the user is logged on the endpoint `logout` because the endpoint doesn't contains sensitive logic.

Now all we have do to is do the handler to delete the cookie.

```go
//Filename: internal/handlers/auth.handlers.go

func (ah *AuthHandler) Logout(c echo.Context) error {
	cookie := http.Cookie{
		Name:     "user",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	c.SetCookie(&cookie)

	return c.Redirect(http.StatusSeeOther, "/")
}
```

And it should works fine.
