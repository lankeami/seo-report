package processor

import (
	"strings"

	"github.com/jaychinthrajah/seo-report/internal/config"
	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

// IsSEORelevant returns true if the item matches any SEO relevance keyword.
func IsSEORelevant(item fetcher.RawItem, keywords []string) bool {
	text := strings.ToLower(item.Title + " " + item.Description)
	for _, kw := range keywords {
		if strings.Contains(text, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// AssignCategory returns the highest-priority category whose keywords match the item.
// Falls back to "Industry News" if no specific category matches.
func AssignCategory(item fetcher.RawItem, categories []config.CategoryConfig) string {
	text := strings.ToLower(item.Title + " " + item.Description)

	bestCategory := "Industry News"
	bestPriority := int(^uint(0) >> 1) // max int

	for _, cat := range categories {
		if cat.Name == "Industry News" {
			continue
		}
		for _, kw := range cat.Keywords {
			if strings.Contains(text, strings.ToLower(kw)) {
				if cat.Priority < bestPriority {
					bestPriority = cat.Priority
					bestCategory = cat.Name
				}
				break
			}
		}
	}

	return bestCategory
}

// Classify filters items for SEO relevance and assigns each a category.
// Returns a map of category name -> items in that category.
func Classify(items []fetcher.RawItem, cfg *config.Config) map[string][]fetcher.RawItem {
	// Collect all SEO relevance keywords
	var relevanceKeywords []string
	for _, kws := range cfg.Keywords {
		relevanceKeywords = append(relevanceKeywords, kws...)
	}

	result := make(map[string][]fetcher.RawItem)
	for _, item := range items {
		if !IsSEORelevant(item, relevanceKeywords) {
			continue
		}
		cat := AssignCategory(item, cfg.Categories)
		result[cat] = append(result[cat], item)
	}
	return result
}
