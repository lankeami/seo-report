package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jaychinthrajah/seo-report/internal/config"
)

func TestFetchAll(t *testing.T) {
	srv1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testFeed))
	}))
	defer srv1.Close()

	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testFeed))
	}))
	defer srv2.Close()

	cfg := &config.Config{
		Sources: config.SourcesConfig{
			RSS: []config.RSSSource{
				{Name: "Source1", URL: srv1.URL, Weight: 3},
				{Name: "Source2", URL: srv2.URL, Weight: 2},
			},
		},
	}

	since := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
	items, errs := FetchAll(cfg, since)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	// Each server returns 1 recent item
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestFetchAll_PartialFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testFeed))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Sources: config.SourcesConfig{
			RSS: []config.RSSSource{
				{Name: "Good", URL: srv.URL, Weight: 1},
				{Name: "Bad", URL: "http://localhost:0/nope", Weight: 1},
			},
		},
	}

	items, errs := FetchAll(cfg, time.Time{})
	if len(errs) == 0 {
		t.Error("expected at least one error for bad URL")
	}
	if len(items) == 0 {
		t.Error("expected items from good source despite bad source")
	}
}
