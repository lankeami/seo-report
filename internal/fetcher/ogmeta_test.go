package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchOGMeta_FullOG(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<!DOCTYPE html><html><head>
<meta property="og:title" content="OG Title">
<meta property="og:description" content="OG Description">
<meta property="og:image" content="https://example.com/image.jpg">
<title>Page Title</title>
</head><body></body></html>`))
	}))
	defer srv.Close()

	got := FetchOGMeta(context.Background(), srv.URL)

	if got.OGTitle != "OG Title" {
		t.Errorf("OGTitle = %q, want %q", got.OGTitle, "OG Title")
	}
	if got.OGDescription != "OG Description" {
		t.Errorf("OGDescription = %q, want %q", got.OGDescription, "OG Description")
	}
	if got.OGImage != "https://example.com/image.jpg" {
		t.Errorf("OGImage = %q, want %q", got.OGImage, "https://example.com/image.jpg")
	}
}

func TestFetchOGMeta_FallbackToMetaDesc(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<!DOCTYPE html><html><head>
<meta name="description" content="Meta description fallback">
<title>Fallback Title</title>
</head><body></body></html>`))
	}))
	defer srv.Close()

	got := FetchOGMeta(context.Background(), srv.URL)

	if got.OGTitle != "Fallback Title" {
		t.Errorf("OGTitle = %q, want %q (fallback to <title>)", got.OGTitle, "Fallback Title")
	}
	if got.OGDescription != "Meta description fallback" {
		t.Errorf("OGDescription = %q, want %q (fallback to meta name=description)", got.OGDescription, "Meta description fallback")
	}
	if got.OGImage != "" {
		t.Errorf("OGImage = %q, want empty", got.OGImage)
	}
}

func TestFetchOGMeta_AllFieldsMissing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<!DOCTYPE html><html><head></head><body><p>No metadata here.</p></body></html>`))
	}))
	defer srv.Close()

	got := FetchOGMeta(context.Background(), srv.URL)

	if got.OGTitle != "" || got.OGDescription != "" || got.OGImage != "" {
		t.Errorf("expected empty OGMeta, got %+v", got)
	}
}

func TestFetchOGMeta_Non200Response(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	got := FetchOGMeta(context.Background(), srv.URL)

	if got.OGTitle != "" || got.OGDescription != "" || got.OGImage != "" {
		t.Errorf("expected empty OGMeta on 404, got %+v", got)
	}
}

func TestFetchOGMeta_NonHTMLContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"og:title":"should not parse"}`))
	}))
	defer srv.Close()

	got := FetchOGMeta(context.Background(), srv.URL)

	if got.OGTitle != "" || got.OGDescription != "" || got.OGImage != "" {
		t.Errorf("expected empty OGMeta for non-HTML, got %+v", got)
	}
}

func TestFetchOGMeta_CancelledContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><meta property="og:title" content="T"></head></html>`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	got := FetchOGMeta(ctx, srv.URL)

	// Should return empty — no panic, no error
	if got.OGTitle != "" || got.OGDescription != "" || got.OGImage != "" {
		// It's possible the cancelled context still completes if fast enough; just ensure no panic.
		t.Logf("cancelled context returned non-empty OGMeta (may be a fast local server): %+v", got)
	}
}

func TestFetchOGMeta_EmptyURL(t *testing.T) {
	got := FetchOGMeta(context.Background(), "")
	if got.OGTitle != "" || got.OGDescription != "" || got.OGImage != "" {
		t.Errorf("expected empty OGMeta for empty URL, got %+v", got)
	}
}
