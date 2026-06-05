package main

import (
	"context"
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
			panic(err)
		}

		items := make([]searchItem, 0, len(*resp.JSON200.Items))
		for _, it := range *resp.JSON200.Items {
			path, _ := it["packagePath"].(string)
			module, _ := it["modulePath"].(string)
			version, _ := it["version"].(string)
			synopsis, _ := it["synopsis"].(string)
			items = append(items, searchItem{Path: path, Module: module, Version: version, Synopsis: synopsis})
		}
		return searchMsg{items}
	}
}

func detailCmd(client *pkgsiteapi.ClientWithResponses, item searchItem) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		path := item.Path
		if item.Module != "" {
			path = item.Module
		}

		limit := 12
		versions, err := client.GetVersionsWithResponse(ctx, path, &pkgsiteapi.GetVersionsParams{Limit: &limit})
		if err != nil {
			panic(err)
		}

		var module, version *string
		if item.Module != "" {
			module = &item.Module
		}
		if item.Version != "" {
			version = &item.Version
		}

		limit = 20
		symbols, err := client.GetSymbolsWithResponse(ctx, item.Path, &pkgsiteapi.GetSymbolsParams{
			Module:  module,
			Version: version,
			Limit:   &limit,
		})
		if err != nil {
			panic(err)
		}

		return detailMsg{path: item.Path, versions: versions.JSON200, symbols: symbols.JSON200}
	}
}
