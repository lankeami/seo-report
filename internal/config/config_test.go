package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Write a minimal config to a temp file
	content := `
sources:
  rss:
    - name: Test Source
      url: https://example.com/feed
      weight: 2
keywords:
  seo_relevance:
    - SEO
categories:
  - name: "Industry News"
    keywords: []
    priority: 9
output:
  dir: docs
dedup:
  levenshtein_threshold: 0.85
`
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(cfg.Sources.RSS) != 1 {
		t.Errorf("expected 1 RSS source, got %d", len(cfg.Sources.RSS))
	}
	if cfg.Sources.RSS[0].Name != "Test Source" {
		t.Errorf("expected name 'Test Source', got %q", cfg.Sources.RSS[0].Name)
	}
	if cfg.Sources.RSS[0].Weight != 2 {
		t.Errorf("expected weight 2, got %d", cfg.Sources.RSS[0].Weight)
	}
	if cfg.Dedup.LevenshteinThreshold != 0.85 {
		t.Errorf("expected threshold 0.85, got %f", cfg.Dedup.LevenshteinThreshold)
	}
}
