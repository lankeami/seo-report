package fetcher

import "time"

// OGMeta holds Open Graph metadata fetched from an article page.
type OGMeta struct {
	OGTitle       string
	OGDescription string
	OGImage       string
}

// RawItem represents a single news item fetched from an RSS feed.
type RawItem struct {
	Title        string
	URL          string
	Description  string
	Source       string
	SourceWeight int
	PublishedAt  time.Time
	OGMeta       OGMeta
}
