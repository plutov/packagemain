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
		if resp == nil || resp.JSON200 == nil {
			return errorMsg{text: "search failed: " + responseError(resp, func(r *pkgsiteapi.GetSearchResponse) *pkgsiteapi.Error { return r.JSONDefault })}
		}

		return searchMsg{resp.JSON200.SearchResults()}
	}
}

func responseError[T any](response *T, defaultError func(*T) *pkgsiteapi.Error) string {
	if response != nil {
		if apiError := defaultError(response); apiError != nil && apiError.Message != nil {
			return *apiError.Message
		}
	}
	return "unexpected response"
}

func detailCmd(client *pkgsiteapi.ClientWithResponses, item *pkgsiteapi.SearchResultData) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		path := item.PackagePath
		if item.ModulePath != "" {
			path = item.ModulePath
		}

		limit := 10
		pseudo := true
		versions, err := client.GetVersionsWithResponse(ctx, path, &pkgsiteapi.GetVersionsParams{
			Limit:  &limit,
			Pseudo: &pseudo,
		})
		if err != nil {
			return errorMsg{text: fmt.Sprintf("loading versions failed: %v", err)}
		}
		if versions == nil || versions.JSON200 == nil {
			return errorMsg{text: "loading versions failed: " + responseError(versions, func(r *pkgsiteapi.GetVersionsResponse) *pkgsiteapi.Error { return r.JSONDefault })}
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
		if symbols == nil || symbols.JSON200 == nil {
			return errorMsg{text: "loading symbols failed: " + responseError(symbols, func(r *pkgsiteapi.GetSymbolsResponse) *pkgsiteapi.Error { return r.JSONDefault })}
		}

		return detailMsg{
			path:     item.PackagePath,
			versions: versions.JSON200,
			symbols:  symbols.JSON200,
		}
	}
}
