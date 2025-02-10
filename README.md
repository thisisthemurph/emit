# Emit

emit is a Go package that helps you build HTTP responses in a clean, fluent, and slightly less rage-inducing way.
Stop fighting `http.ResponseWriter` and let emit do the dirty work for you.

## Features

- **Fluent API:** Chain methods like a civilized person instead of sprinkling w.WriteHeader calls everywhere.
- **Convenient Response Methods:** Send text, JSON, error responses, or a 204 No Content (for when you really have nothing to say).
- **Encapsulation:** Keeps your response logic neat, so your teammates don’t hate you.
- **Sensible defaults:** Defaults to status code 200 when appropriate and sets appropriate `Content-Type` headers, because let's be honest, you always forget to set them.

## Instillation

```shell
go get github.com/thisisthemurph/emit
```

## Usage

**Basic usage**

Sometimes you just want to send some text and move on.
- sets the `Content-Type` header to `text/plain`
- applies the status code, with a default of 200
- writes the text, nothing special

```go
package main

import (
	"net/http"
	"github.com/thisisthemurph/emit"
)

func handler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).Text("Hello, world!")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

**Sending a JSON response**

JSON encoding without manually writing `if err := json.NewEncoder(rb.w).Encode(data); err != nil` every time, like a caveman.
- applies the Content-Type application/json header
- writes the status code with a 200 default if you didn't specify one
- encodes the data to JSON

```go
package main

import (
	"net/http"
	"github.com/thisisthemurph/emit"
)

type Response struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).JSON(Response{Message: "Hello, JSON!"})
}
```

**Setting the response status code**

Bored with sending 200 OK status codes?

```go
package main

import (
	"net/http"
	"github.com/thisisthemurph/emit"
)

type Product struct {
	Name string `json:"name"`
	Quantity int `json:"quantity"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).Status(http.StatusCreated).JSON(Product{
		Name: "Henry Hoover",
		Quantity: 1000,
	})
}
```

**Sending no content**

You say it best, when you say nothing at all!

```go
package main

import (
	"net/http"
	"github.com/thisisthemurph/emit"
)

func handler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).NoContent()
}
```

**Setting headers and cookies**

Want to sprinkle in some headers and cookies?

```go
package main

import (
	"net/http"
	"github.com/thisisthemurph/emit"
)

func handler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).
		Header("X-Custom-Header", "Value").
		Cookie(&http.Cookie{Name: "session", Value: "abc123"}).
		Text("Headers and cookies set!")
}
```

**Handling errors**

The `ErrorJSON` method is an *opinionated* (and very limited) way to send an error response.
- Uses the JSON format: `{"message": "Something went wrong"}`
- Defaults to 500 Internal Server Error, because what else could it be?

```go
package main

import (
	"net/http"
	"github.com/thisisthemurph/emit"
)

func handler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).ErrorJSON("Something went wrong")
}

func anotherHandler(w http.ResponseWriter, r *http.Request) {
	emit.New(w).Status(http.StatusBadRequest).ErrorJSON("You didn't star the repo!")
}
```
