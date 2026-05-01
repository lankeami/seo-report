package fetcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jaychinthrajah/seo-report/internal/config"
)

// FetchAll fetches all RSS sources in parallel and returns all items and any errors.
func FetchAll(cfg *config.Config, since time.Time) ([]RawItem, []error) {
	var mu sync.Mutex
	var allItems []RawItem
	var allErrors []error
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, src := range cfg.Sources.RSS {
		wg.Add(1)
		go func(s config.RSSSource) {
			defer wg.Done()
			items, err := fetchWithContext(ctx, s, since)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				allErrors = append(allErrors, fmt.Errorf("source %q: %w", s.Name, err))
				return
			}
			allItems = append(allItems, items...)
		}(src)
	}

	wg.Wait()
	return allItems, allErrors
}

func fetchWithContext(ctx context.Context, src config.RSSSource, since time.Time) ([]RawItem, error) {
	done := make(chan struct{})
	var items []RawItem
	var fetchErr error

	go func() {
		items, fetchErr = FetchRSS(src.URL, src.Name, src.Weight, since)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout fetching feed: %w", ctx.Err())
	case <-done:
		return items, fetchErr
	}
}
