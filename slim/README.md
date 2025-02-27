# Slim

Slim is a template engine created by [mattn](https://github.com/mattn/go-slim), to see the original syntax documentation please [click here](https://rubydoc.info/gems/slim/frames)

### Basic Example

_**./views/index.slim**_
```html
== render("partials/header.slim")

h1 = Title

== render("partials/footer.slim")
```
_**./views/partials/header.slim**_
```html
h2 = Header
```
_**./views/partials/footer.slim**_
```html
h2 = Footer
```
_**./views/layouts/main.slim**_
```html
doctype html
html
  head
    title Main
    include ../partials/meta.slim
  body
    | {{embed}}
```

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/slim/v2"

	// "net/http" // embedded system
)

func main() {
	// Create a new engine
	engine := slim.New("./views", ".slim")

	// Or from an embedded system
	// See github.com/gofiber/embed for examples
	// engine := slim.NewFileSystem(http.Dir("./views", ".slim"))

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/layout", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}

```
