package main

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/plutov/gopkg/pkgsiteapi"
)

func searchCmd(client *pkgsiteapi.ClientWithResponses, q string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		limit := 10
		resp, err := client.GetSearchWithResponse(ctx, &pkgsiteapi.GetSearchParams{Q: &q, Limit: &limit})
		if err != nil {
			return errorMsg{text: fmt.Sprintf("search failed: %v", err)}
		}

		return searchMsg{resp.JSON200.SearchResults()}
	}
}

func detailCmd(client *pkgsiteapi.ClientWithResponses, item *pkgsiteapi.SearchResult) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		path := item.PackagePath
		if item.ModulePath != "" {
			path = item.ModulePath
		}

		limit := 10
		versions, err := client.GetVersionsWithResponse(ctx, path, &pkgsiteapi.GetVersionsParams{Limit: &limit})
		if err != nil {
			return errorMsg{text: fmt.Sprintf("loading versions failed: %v", err)}
		}

		var module, version *string
		if item.ModulePath != "" {
			module = &item.ModulePath
		}
		if item.Version != "" {
			version = &item.Version
		}

		limit = 10
		symbols, err := client.GetSymbolsWithResponse(ctx, item.PackagePath, &pkgsiteapi.GetSymbolsParams{
			Module:  module,
			Version: version,
			Limit:   &limit,
		})
		if err != nil {
			return errorMsg{text: fmt.Sprintf("loading symbols failed: %v", err)}
		}

		return detailMsg{
			path:     item.PackagePath,
			versions: versions.JSON200,
			symbols:  symbols.JSON200,
		}
	}
}
