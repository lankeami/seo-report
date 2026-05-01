package processor

import (
	"strings"

	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

var urgencyKeywords = []string{"breaking", "just announced", "critical", "major update"}

// Score computes an impact score 1–100 for a classified item.
// base: 50
// + source weight bonus: weight 1→+0, 2→+10, 3→+20
// + coverage: each extra source covering same story +5, capped at +20
// + urgency: title contains urgency keyword → +10
// + AEO bonus: category is "AI Overviews & AEO" → +10
func Score(item fetcher.RawItem, category string, coverageCount int) int {
	score := 50

	// Source weight bonus
	switch item.SourceWeight {
	case 2:
		score += 10
	case 3:
		score += 20
	}

	// Coverage bonus (capped at 20)
	coverage := (coverageCount - 1) * 5
	if coverage > 20 {
		coverage = 20
	}
	if coverage > 0 {
		score += coverage
	}

	// Urgency bonus
	titleLower := strings.ToLower(item.Title)
	for _, kw := range urgencyKeywords {
		if strings.Contains(titleLower, kw) {
			score += 10
			break
		}
	}

	// AEO bonus
	if category == "AI Overviews & AEO" {
		score += 10
	}

	// Clamp to 1–100
	if score < 1 {
		score = 1
	}
	if score > 100 {
		score = 100
	}

	return score
}
