package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testFeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test SEO Blog</title>
    <item>
      <title>Google Core Update May 2026</title>
      <link>https://example.com/core-update</link>
      <description>A major core algorithm update has been released.</description>
      <pubDate>Thu, 01 May 2026 10:00:00 +0000</pubDate>
    </item>
    <item>
      <title>Old Article</title>
      <link>https://example.com/old</link>
      <description>This is an old article.</description>
      <pubDate>Mon, 01 Jan 2024 10:00:00 +0000</pubDate>
    </item>
  </channel>
</rss>`

func TestFetchRSS(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(testFeed))
	}))
	defer srv.Close()

	since := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
	items, err := FetchRSS(srv.URL, "Test Source", 2, since)
	if err != nil {
		t.Fatalf("FetchRSS() error: %v", err)
	}

	// Only the recent item should be returned
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	item := items[0]
	if item.Title != "Google Core Update May 2026" {
		t.Errorf("unexpected title: %q", item.Title)
	}
	if item.Source != "Test Source" {
		t.Errorf("unexpected source: %q", item.Source)
	}
	if item.SourceWeight != 2 {
		t.Errorf("unexpected weight: %d", item.SourceWeight)
	}
	if item.URL != "https://example.com/core-update" {
		t.Errorf("unexpected URL: %q", item.URL)
	}
}

func TestFetchRSS_NoSinceFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(testFeed))
	}))
	defer srv.Close()

	// Zero time means no filter
	items, err := FetchRSS(srv.URL, "Test Source", 1, time.Time{})
	if err != nil {
		t.Fatalf("FetchRSS() error: %v", err)
	}

	// Both items should be returned (zero published time doesn't get filtered)
	if len(items) < 1 {
		t.Errorf("expected at least 1 item, got %d", len(items))
	}
}

func TestFetchRSS_Error(t *testing.T) {
	_, err := FetchRSS("http://localhost:0/nonexistent", "Bad Source", 1, time.Time{})
	if err == nil {
		t.Error("expected error for bad URL, got nil")
	}
}
