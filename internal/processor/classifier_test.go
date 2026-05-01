package processor

import (
	"testing"
	"time"

	"github.com/jaychinthrajah/seo-report/internal/config"
	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

var testCategories = []config.CategoryConfig{
	{Name: "AI Overviews & AEO", Keywords: []string{"ai overview", "answer engine", "aeo"}, Priority: 1},
	{Name: "Algorithm Updates", Keywords: []string{"core update", "algorithm update", "google update"}, Priority: 2},
	{Name: "Technical SEO", Keywords: []string{"crawling", "indexing", "core web vitals"}, Priority: 3},
	{Name: "Industry News", Keywords: []string{}, Priority: 9},
}

var testKeywords = map[string][]string{
	"seo_relevance": {"SEO", "google search", "serp", "algorithm", "ai overview", "aeo", "crawling"},
}

func makeItem(title, desc string) fetcher.RawItem {
	return fetcher.RawItem{Title: title, Description: desc, PublishedAt: time.Now()}
}

func TestIsSEORelevant(t *testing.T) {
	keywords := []string{"SEO", "google search", "algorithm"}

	tests := []struct {
		title string
		desc  string
		want  bool
	}{
		{"Google SEO Update", "", true},
		{"Baking Bread at Home", "flour and water", false},
		{"Tech News", "google search ranking change", true},
	}

	for _, tt := range tests {
		item := makeItem(tt.title, tt.desc)
		if got := IsSEORelevant(item, keywords); got != tt.want {
			t.Errorf("IsSEORelevant(%q) = %v, want %v", tt.title, got, tt.want)
		}
	}
}

func TestAssignCategory(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"Google AI Overview Changes Everything", "AI Overviews & AEO"},
		{"Google Core Update Released Today", "Algorithm Updates"},
		{"Improving Site Crawling Speed", "Technical SEO"},
		{"General SEO News This Week", "Industry News"},
	}

	for _, tt := range tests {
		item := makeItem(tt.title, "")
		got := AssignCategory(item, testCategories)
		if got != tt.want {
			t.Errorf("AssignCategory(%q) = %q, want %q", tt.title, got, tt.want)
		}
	}
}

func TestClassify(t *testing.T) {
	cfg := &config.Config{
		Keywords:   testKeywords,
		Categories: testCategories,
	}

	items := []fetcher.RawItem{
		makeItem("Google Core Update", "algorithm update released"),
		makeItem("Baking Tips", "not related at all"),
		makeItem("AEO Best Practices", "answer engine optimization guide"),
	}

	classified := Classify(items, cfg)

	if _, ok := classified["Baking Tips"]; ok {
		t.Error("non-SEO item should not be classified")
	}

	if len(classified["Algorithm Updates"]) != 1 {
		t.Errorf("expected 1 Algorithm Updates item, got %d", len(classified["Algorithm Updates"]))
	}

	if len(classified["AI Overviews & AEO"]) != 1 {
		t.Errorf("expected 1 AEO item, got %d", len(classified["AI Overviews & AEO"]))
	}
}
