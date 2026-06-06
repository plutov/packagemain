package pkgsiteapi

import "time"

type SearchResult struct {
	PackagePath string `json:"packagePath,omitempty"`
	ModulePath  string `json:"modulePath,omitempty"`
	Version     string `json:"version,omitempty"`
	Synopsis    string `json:"synopsis,omitempty"`
}

type VersionResult struct {
	Version    string     `json:"version,omitempty"`
	CommitTime *time.Time `json:"commitTime,omitempty"`
}

type SymbolResult struct {
	Kind     string `json:"kind,omitempty"`
	Name     string `json:"name,omitempty"`
	Synopsis string `json:"synopsis,omitempty"`
}

func (p *PaginatedResponse) SearchResults() []SearchResult {
	if p == nil || p.Items == nil {
		return nil
	}

	items := make([]SearchResult, 0, len(*p.Items))
	for _, item := range *p.Items {
		items = append(items, SearchResult{
			PackagePath: stringValue(item, "packagePath"),
			ModulePath:  stringValue(item, "modulePath"),
			Version:     stringValue(item, "version"),
			Synopsis:    stringValue(item, "synopsis"),
		})
	}

	return items
}

func (p *PaginatedResponse) VersionResults() []VersionResult {
	if p == nil || p.Items == nil {
		return nil
	}

	items := make([]VersionResult, 0, len(*p.Items))
	for _, item := range *p.Items {
		items = append(items, VersionResult{
			Version:    stringValue(item, "version"),
			CommitTime: timeValue(item, "commitTime"),
		})
	}

	return items
}

func (p *PaginatedResponse) SymbolResults() []SymbolResult {
	if p == nil || p.Items == nil {
		return nil
	}

	items := make([]SymbolResult, 0, len(*p.Items))
	for _, item := range *p.Items {
		items = append(items, SymbolResult{
			Kind:     stringValue(item, "kind"),
			Name:     stringValue(item, "name"),
			Synopsis: stringValue(item, "synopsis"),
		})
	}

	return items
}

func (p *PackageSymbols) SymbolResults() []SymbolResult {
	if p == nil || p.Symbols == nil {
		return nil
	}
	return p.Symbols.SymbolResults()
}

func stringValue(item map[string]interface{}, key string) string {
	value, _ := item[key].(string)
	return value
}

func timeValue(item map[string]interface{}, key string) *time.Time {
	value, _ := item[key].(string)
	if value == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}

	return &t
}
