package fetcher

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

// FetchRSS fetches items from an RSS/Atom feed published after `since`.
func FetchRSS(url, sourceName string, sourceWeight int, since time.Time) ([]RawItem, error) {
	parser := gofeed.NewParser()
	parser.UserAgent = "seo-report/1.0 (github.com/jaychinthrajah/seo-report)"
	feed, err := parser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parsing feed %s: %w", sourceName, err)
	}

	var items []RawItem
	for _, entry := range feed.Items {
		published := time.Time{}
		if entry.PublishedParsed != nil {
			published = *entry.PublishedParsed
		} else if entry.UpdatedParsed != nil {
			published = *entry.UpdatedParsed
		}

		if !published.IsZero() && published.Before(since) {
			continue
		}

		items = append(items, RawItem{
			Title:        entry.Title,
			URL:          entry.Link,
			Description:  entry.Description,
			Source:       sourceName,
			SourceWeight: sourceWeight,
			PublishedAt:  published,
		})
	}

	return items, nil
}
