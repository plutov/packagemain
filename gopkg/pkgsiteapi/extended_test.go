package pkgsiteapi

import (
	"testing"
	"time"
)

func TestTypedPaginationAdapters(t *testing.T) {
	commitTime := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	packagePath := "example.com/module/package"
	version := "v1.2.3"
	kind := "Function"
	name := "Do"

	search := (&PaginatedResponseSearchResult{Items: &[]SearchResult{{
		PackagePath: &packagePath,
	}}}).SearchResults()
	if len(search) != 1 || search[0].PackagePath != packagePath || search[0].ModulePath != "" {
		t.Fatalf("SearchResults() = %#v", search)
	}

	versions := (&PaginatedResponseModuleVersion{Items: &[]ModuleVersion{{
		Version:    &version,
		CommitTime: &commitTime,
	}}}).VersionResults()
	if len(versions) != 1 || versions[0].Version != version || !versions[0].CommitTime.Equal(commitTime) {
		t.Fatalf("VersionResults() = %#v", versions)
	}

	symbols := (&PackageSymbols{Symbols: &PaginatedResponseSymbol{Items: &[]Symbol{{
		Kind: &kind,
		Name: &name,
	}}}}).SymbolResults()
	if len(symbols) != 1 || symbols[0].Kind != kind || symbols[0].Name != name || symbols[0].Synopsis != "" {
		t.Fatalf("SymbolResults() = %#v", symbols)
	}
}

func TestTypedPaginationAdaptersHandleNil(t *testing.T) {
	if results := (*PaginatedResponseSearchResult)(nil).SearchResults(); results != nil {
		t.Fatalf("nil search page results = %#v", results)
	}
	if results := (&PaginatedResponseModuleVersion{}).VersionResults(); results != nil {
		t.Fatalf("empty version page results = %#v", results)
	}
	if results := (*PackageSymbols)(nil).SymbolResults(); results != nil {
		t.Fatalf("nil package symbols results = %#v", results)
	}
}
