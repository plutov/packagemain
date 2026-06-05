For years, pkg.go.dev has been the central place to discover Go packages. It's a great project, I like the consistent clean design. Every package is there, it doesn't have to be explicitly published.

For example I am maintaining this Go library, and it's already there.
https://pkg.go.dev/github.com/plutov/paypal/v4

And the amazing Go team recently has launched the pkg.go.dev API, giving us the direct access to package metadata, versions, symbols, vulnerabilities, import relationships, and more.

https://go.dev/blog/pkgsite-api

The machine-readable API contract is also published directly as an which will be useful for us in a second.

https://pkg.go.dev/v1beta/openapi.yaml

Since I am working mostly in the terminal using Neovim, I can now browse the Go docs not leaving it, which is a huge performance booster.

So, in this video we're going to create a TUI explorer that lets us:

- Search for packages
- Browse module versions
- Inspect package symbols
- View import relationships
- Check vulnerabilities

All directly from the new pkg.go.dev API.

Let's get started.

## What Is the pkg.go.dev API?

It exposes structured metadata from pkg.go.dev through a collection of REST endpoints.

Some of the most useful endpoints include:

- /search
- /packagу
- /module
- /versions
- /symbols
- /imported-by
- /vulns
    
Everything is GET-only and designed to be cache-friendly.

For example, if we want information about the cmp package:

```bash
curl https://pkg.go.dev/v1beta/package/github.com/google/go-cmp/cmp | jq
```

The response includes:

- module path
- version
- package name
- synopsis
- redistributable status

This means we can build tools without scraping a single HTML page.

## TUI Concept

Think of it as a lightweight package browser directly inside your terminal.


## API Client

Since there is a OpenAPI schema, we can generate a client SDK easily, using the nice package from Jamie.

https://github.com/oapi-codegen/oapi-codegen

I already had a video about oapi-codegen: LINK

```bash
mkdir gopkg/pkgsiteapi
cd gopkg

go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

oapi-codegen \
  -generate types,client \
  -package pkgsiteapi \
  https://pkg.go.dev/v1beta/openapi.yaml > pkgsiteapi/client.gen.go
```

Once we have the spec locally and the client generated, we can move on to the terminal UI itself.

## Scaffolding the TUI

Now that we have the client ready, we can start building the terminal UI.

I want to keep the first version very small: one `main.go` file, one Bubble Tea model, one text input, and one list of results. No extra packages yet, no complex navigation, and no premature architecture. Goal is to get something working quickly and then grow it step by step.

First, install Bubble Tea and a couple of useful Bubbles components:

```bash
go mod init github.com/plutov/gopkg
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles/textinput
```

At this point, we can create a `main.go` file and define a single model for the whole app.

For the first pass, that model only needs a few fields:

- a `textinput.Model` for the search query
- a slice of search results from the generated SDK
- an integer for the selected row
- a `loading` flag
- an `err` field

That gives us enough state to support the core interaction:

1. type a package query
2. press Enter
3. load results from pkg.go.dev
4. move through them with the keyboard

This is a good fit for Bubble Tea because the program naturally breaks down into the usual three parts:

- `Init` for startup behavior
- `Update` for handling keys and async responses
- `View` for rendering the search box and result list

In the very first version, I would not use multiple screens yet. I would keep it to a single view: input on top, results underneath, and a tiny help line at the bottom. Once that feels good, we can add a details pane for the selected package.

So next we are going to write a minimal `main.go` that boots Bubble Tea, focuses the input, and renders an empty search screen. After that, we can hook Enter up to the generated client and make the first real API request.

Btw, there is already a video on this channel that covers the basics of Bubble Tea, so if you want to follow along with the code, check it out: LINK
