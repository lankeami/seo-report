package renderer

import (
	"strings"
	"testing"
	"time"
)

func TestRenderIndex_ContainsKeyElements(t *testing.T) {
	reports := []ReportMeta{
		{
			Date:          time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
			Filename:      "2026-05-01.html",
			ItemCount:     12,
			CategoryCount: 4,
		},
		{
			Date:          time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
			Filename:      "2026-04-30.html",
			ItemCount:     8,
			CategoryCount: 3,
		},
	}

	html := RenderIndex(reports, time.Now())

	checks := []string{
		"AEO &amp; SEO Daily Digest",
		"Daily aggregation of AEO and SEO news",
		"2026-05-01.html",
		"2026-04-30.html",
		`target="_blank"`,
		`rel="noopener noreferrer"`,
		"12 items",
		"4 categories",
		"Last updated:",
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("index HTML missing expected content: %q", check)
		}
	}
}

func TestRenderIndex_Empty(t *testing.T) {
	html := RenderIndex(nil, time.Now())
	if !strings.Contains(html, "No reports yet") {
		t.Error("expected empty state message when no reports")
	}
}
