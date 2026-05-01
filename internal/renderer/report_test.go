package renderer

import (
	"strings"
	"testing"
	"time"

	"github.com/jaychinthrajah/seo-report/internal/config"
	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

func TestBuildReport(t *testing.T) {
	categories := []config.CategoryConfig{
		{Name: "Algorithm Updates", Keywords: []string{"core update"}, Priority: 2},
		{Name: "Industry News", Keywords: []string{}, Priority: 9},
	}

	classified := map[string][]fetcher.RawItem{
		"Algorithm Updates": {
			{Title: "Google Core Update", URL: "https://example.com/1", Source: "SEJ", SourceWeight: 3, PublishedAt: time.Now()},
		},
	}

	result := BuildReport(classified, categories)
	if len(result) != 1 {
		t.Fatalf("expected 1 category, got %d", len(result))
	}
	if result[0].Name != "Algorithm Updates" {
		t.Errorf("expected 'Algorithm Updates', got %q", result[0].Name)
	}
	if len(result[0].Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(result[0].Items))
	}
}

func TestRenderReport_ContainsKeyElements(t *testing.T) {
	date := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	generatedAt := time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC)

	cats := []CategoryItems{
		{
			Name:     "Algorithm Updates",
			Priority: 2,
			Items: []ReportItem{
				{
					Title:       "Google Core Update",
					URL:         "https://example.com/update",
					Source:      "Search Engine Journal",
					PublishedAt: date,
					Score:       75,
				},
			},
		},
	}

	html := RenderReport(date, cats, []string{"Search Engine Journal"}, generatedAt)

	checks := []string{
		"AEO &amp; SEO Daily Digest",
		"2026-05-01",
		"Algorithm Updates",
		"Google Core Update",
		`href="https://example.com/update"`,
		`target="_blank"`,
		`rel="noopener noreferrer"`,
		"Score: 75",
		"Search Engine Journal",
		"Generated:",
		`href="index.html"`,
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("rendered HTML missing expected content: %q", check)
		}
	}
}

func TestRenderReport_EmptyCategories(t *testing.T) {
	date := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	html := RenderReport(date, nil, nil, time.Now())
	if !strings.Contains(html, "No SEO/AEO items found") {
		t.Error("expected empty state message when no categories")
	}
}
