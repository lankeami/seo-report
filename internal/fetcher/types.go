package fetcher

import "time"

// RawItem represents a single news item fetched from an RSS feed.
type RawItem struct {
	Title       string
	URL         string
	Description string
	Source      string
	SourceWeight int
	PublishedAt time.Time
}
