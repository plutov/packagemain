package pkgsiteapi

import "time"

// SearchResultData is the UI-friendly representation of a search result.
type SearchResultData struct {
	PackagePath string
	ModulePath  string
	Version     string
	Synopsis    string
}

type VersionResult struct {
	Version    string
	CommitTime *time.Time
}

type SymbolResult struct {
	Kind     string
	Name     string
	Synopsis string
}

func (p *PaginatedResponseSearchResult) SearchResults() []SearchResultData {
	if p == nil || p.Items == nil {
		return nil
	}

	items := make([]SearchResultData, 0, len(*p.Items))
	for _, item := range *p.Items {
		items = append(items, SearchResultData{
			PackagePath: derefString(item.PackagePath),
			ModulePath:  derefString(item.ModulePath),
			Version:     derefString(item.Version),
			Synopsis:    derefString(item.Synopsis),
		})
	}
	return items
}

func (p *PaginatedResponseModuleVersion) VersionResults() []VersionResult {
	if p == nil || p.Items == nil {
		return nil
	}

	items := make([]VersionResult, 0, len(*p.Items))
	for _, item := range *p.Items {
		items = append(items, VersionResult{
			Version:    derefString(item.Version),
			CommitTime: item.CommitTime,
		})
	}
	return items
}

func (p *PaginatedResponseSymbol) SymbolResults() []SymbolResult {
	if p == nil || p.Items == nil {
		return nil
	}

	items := make([]SymbolResult, 0, len(*p.Items))
	for _, item := range *p.Items {
		items = append(items, SymbolResult{
			Kind:     derefString(item.Kind),
			Name:     derefString(item.Name),
			Synopsis: derefString(item.Synopsis),
		})
	}
	return items
}

func (p *PackageSymbols) SymbolResults() []SymbolResult {
	if p == nil {
		return nil
	}
	return p.Symbols.SymbolResults()
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
