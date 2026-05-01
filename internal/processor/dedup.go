package processor

import (
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

// ScoredItem is a classified, deduplicated, and scored news item.
type ScoredItem struct {
	Title         string
	URL           string
	Description   string
	Source        string
	SourceWeight  int
	PublishedAt   interface{} // time.Time, kept as interface to avoid import cycle
	Category      string
	Score         int
	CoverageCount int            // number of sources covering this story
	ExtraSources  []string       // additional source names (beyond primary)
}

// DeduplicatedItem groups items covering the same story.
type DeduplicatedItem struct {
	fetcher.RawItem
	Sources []string // all source names covering this story
}

func titleSimilarity(a, b string) float64 {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	if a == b {
		return 1.0
	}

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	if maxLen == 0 {
		return 1.0
	}

	dist := levenshtein.ComputeDistance(a, b)
	return 1.0 - float64(dist)/float64(maxLen)
}

// Deduplicate merges items with identical URLs or similar titles (above threshold).
// Returns deduplicated items, each tracking all sources that covered the story.
func Deduplicate(items []fetcher.RawItem, threshold float64) []DeduplicatedItem {
	var result []DeduplicatedItem

	for _, item := range items {
		merged := false
		for i, existing := range result {
			sameURL := item.URL != "" && item.URL == existing.URL
			similarTitle := titleSimilarity(item.Title, existing.Title) >= threshold
			if sameURL || similarTitle {
				result[i].Sources = append(result[i].Sources, item.Source)
				merged = true
				break
			}
		}

		if !merged {
			result = append(result, DeduplicatedItem{
				RawItem: item,
				Sources: []string{item.Source},
			})
		}
	}

	return result
}
