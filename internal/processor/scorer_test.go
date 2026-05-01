package processor

import (
	"testing"

	"github.com/jaychinthrajah/seo-report/internal/fetcher"
)

func TestScore(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		weight        int
		category      string
		coverageCount int
		wantMin       int
		wantMax       int
	}{
		{
			name:          "base score weight 1",
			title:         "Some SEO news",
			weight:        1,
			category:      "Industry News",
			coverageCount: 1,
			wantMin:       50,
			wantMax:       50,
		},
		{
			name:          "weight 3 source",
			title:         "Some SEO news",
			weight:        3,
			category:      "Industry News",
			coverageCount: 1,
			wantMin:       70,
			wantMax:       70,
		},
		{
			name:          "AEO category bonus",
			title:         "AI Overview changes",
			weight:        1,
			category:      "AI Overviews & AEO",
			coverageCount: 1,
			wantMin:       60,
			wantMax:       60,
		},
		{
			name:          "urgency keyword",
			title:         "Breaking: Google Update",
			weight:        1,
			category:      "Industry News",
			coverageCount: 1,
			wantMin:       60,
			wantMax:       60,
		},
		{
			name:          "coverage bonus capped",
			title:         "Covered everywhere",
			weight:        1,
			category:      "Industry News",
			coverageCount: 10,
			wantMin:       70,
			wantMax:       70,
		},
		{
			name:          "all bonuses capped at 100",
			title:         "Breaking: Major Update Just Announced",
			weight:        3,
			category:      "AI Overviews & AEO",
			coverageCount: 10,
			wantMin:       100,
			wantMax:       100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := fetcher.RawItem{Title: tt.title, SourceWeight: tt.weight}
			got := Score(item, tt.category, tt.coverageCount)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("Score() = %d, want [%d, %d]", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}
