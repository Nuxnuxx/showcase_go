# Scaffolding

```bash
go mod init github.com/user/project
```

After that create a ```main.go``` and add it 
```go
// Filename: main.go
package main

func main(){

}
```

This is the entrypoint of your binary

# Setting up

### Intall echo

```bash
go get github.com/labstack/echo/v4
```

### Get your Api Key for RAWG

What is RAWG : [Here](https://rawg.io/)

[Get your api key here](https://rawg.io/apidocs)

### Load .env

This package his self explanatory, it just autoload what is in the ```.env``` automatically for us

```bash
go get "github.com/jpfuentes2/go-env/autoload"
```

```go
// Filename: main.go
import (
    ...
    _ "github.com/jpfuentes2/go-env/autoload"
)
```

### Get started with RAWG

```go
// Filename: main.go

// GET API_KEY from .env file
API_KEY := os.Getenv("API_KEY")
// Remove quotes from API_KEY
API_KEY = strings.Trim(API_KEY, `"`)

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

### Lets run this

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

```make run``` it will build and then start the binary

### ---- Checkpoint
```git checkout 779bf031aa778ac5dbb42194b4cdbabee8fc93c5```
### ---- Checkpoint
