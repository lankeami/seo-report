package processor

import (
	"testing"

	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

func rawItem(title, url, source string) fetcher.RawItem {
	return fetcher.RawItem{Title: title, URL: url, Source: source}
}

func TestDeduplicate_SameURL(t *testing.T) {
	items := []fetcher.RawItem{
		rawItem("Google Update", "https://example.com/update", "SEJ"),
		rawItem("Google Update News", "https://example.com/update", "SEL"),
	}

	result := Deduplicate(items, 0.85)
	if len(result) != 1 {
		t.Fatalf("expected 1 deduplicated item, got %d", len(result))
	}
	if len(result[0].Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(result[0].Sources))
	}
}

func TestDeduplicate_SimilarTitles(t *testing.T) {
	items := []fetcher.RawItem{
		rawItem("Google Core Algorithm Update 2026", "https://a.com/1", "Source1"),
		rawItem("Google Core Algorithm Update 2026!", "https://b.com/2", "Source2"),
	}

	result := Deduplicate(items, 0.85)
	if len(result) != 1 {
		t.Fatalf("expected 1 deduplicated item (similar titles), got %d", len(result))
	}
}

func TestDeduplicate_DifferentItems(t *testing.T) {
	items := []fetcher.RawItem{
		rawItem("Google Core Update", "https://a.com/1", "S1"),
		rawItem("Ahrefs Releases New Keyword Tool", "https://b.com/2", "S2"),
		rawItem("Local SEO Tips for 2026", "https://c.com/3", "S3"),
	}

	result := Deduplicate(items, 0.85)
	if len(result) != 3 {
		t.Errorf("expected 3 items (all distinct), got %d", len(result))
	}
}

func TestTitleSimilarity(t *testing.T) {
	tests := []struct {
		a, b    string
		wantMin float64
	}{
		{"identical title", "identical title", 1.0},
		{"google update", "google update!", 0.9},
		{"completely different", "nothing alike here", 0.0},
	}

	for _, tt := range tests {
		sim := titleSimilarity(tt.a, tt.b)
		if sim < tt.wantMin {
			t.Errorf("titleSimilarity(%q, %q) = %.2f, want >= %.2f", tt.a, tt.b, sim, tt.wantMin)
		}
	}
}
