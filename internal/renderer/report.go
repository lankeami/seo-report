package renderer

import (
	"fmt"
	"html"
	"sort"
	"strings"
	"time"

	"github.com/jaychinthrajah/seo-report/internal/config"
	"github.com/jaychinthrajah/seo-report/internal/fetcher"
	"github.com/jaychinthrajah/seo-report/internal/processor"
)

// CategoryItems holds scored, deduplicated items for a category.
type CategoryItems struct {
	Name     string
	Priority int
	Items    []ReportItem
}

// ReportItem is a fully processed item ready to render.
type ReportItem struct {
	Title         string
	URL           string
	Source        string
	ExtraSources  []string
	PublishedAt   time.Time
	Score         int
	CoverageCount int
}

// BuildReport processes classified items into sorted, scored report data.
func BuildReport(classified map[string][]fetcher.RawItem, categories []config.CategoryConfig) []CategoryItems {
	// Build priority lookup
	priorityMap := make(map[string]int, len(categories))
	for _, cat := range categories {
		priorityMap[cat.Name] = cat.Priority
	}

	var result []CategoryItems

	for catName, items := range classified {
		if len(items) == 0 {
			continue
		}

		// Deduplicate within category
		deduped := processor.Deduplicate(items, 0.85)

		// Score each item
		var reportItems []ReportItem
		for _, d := range deduped {
			score := processor.Score(d.RawItem, catName, len(d.Sources))
			extra := d.Sources[1:]
			reportItems = append(reportItems, ReportItem{
				Title:         d.Title,
				URL:           d.URL,
				Source:        d.Source,
				ExtraSources:  extra,
				PublishedAt:   d.PublishedAt,
				Score:         score,
				CoverageCount: len(d.Sources),
			})
		}

		// Sort by score descending
		sort.Slice(reportItems, func(i, j int) bool {
			return reportItems[i].Score > reportItems[j].Score
		})

		priority := 99
		if p, ok := priorityMap[catName]; ok {
			priority = p
		}

		result = append(result, CategoryItems{
			Name:     catName,
			Priority: priority,
			Items:    reportItems,
		})
	}

	// Sort categories by priority
	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

// RenderReport renders the per-day HTML report.
func RenderReport(date time.Time, categories []CategoryItems, sourceNames []string, generatedAt time.Time) string {
	dateStr := date.Format("2006-01-02")
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AEO &amp; SEO Daily Digest &mdash; ` + html.EscapeString(dateStr) + `</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:'Inter',sans-serif;background:#f8fafc;color:#1e293b;line-height:1.6;padding:0 16px}
.container{max-width:900px;margin:0 auto;padding:32px 0}
header{margin-bottom:32px;padding-bottom:24px;border-bottom:2px solid #0d9488}
h1{font-size:1.75rem;font-weight:600;color:#0d9488}
.subtitle{color:#64748b;margin-top:4px;font-size:0.95rem}
.back-link{display:inline-block;margin-bottom:16px;color:#0d9488;text-decoration:none;font-size:0.9rem}
.back-link:hover{text-decoration:underline}
details{background:#fff;border:1px solid #e2e8f0;border-radius:8px;margin-bottom:16px;overflow:hidden}
summary{padding:16px 20px;cursor:pointer;font-weight:600;font-size:1rem;color:#1e293b;display:flex;align-items:center;gap:8px;list-style:none}
summary::-webkit-details-marker{display:none}
summary::before{content:'▶';color:#0d9488;font-size:0.75rem;transition:transform 0.2s}
details[open] summary::before{transform:rotate(90deg)}
.item-count{background:#f1f5f9;color:#64748b;font-size:0.75rem;font-weight:500;padding:2px 8px;border-radius:12px;margin-left:auto}
.items{padding:0 20px 16px}
.item{padding:12px 0;border-bottom:1px solid #f1f5f9}
.item:last-child{border-bottom:none}
.item-title a{color:#1e293b;text-decoration:none;font-weight:500;font-size:0.95rem}
.item-title a:hover{color:#0d9488}
.item-meta{display:flex;flex-wrap:wrap;gap:6px;margin-top:6px;align-items:center}
.badge{display:inline-block;padding:2px 8px;border-radius:12px;font-size:0.75rem;font-weight:500}
.badge-source{background:#dbeafe;color:#1e40af}
.badge-score{background:#dcfce7;color:#15803d}
.badge-score-high{background:#d1fae5;color:#065f46}
.badge-time{color:#94a3b8;font-size:0.75rem}
footer{margin-top:48px;padding-top:24px;border-top:1px solid #e2e8f0;color:#94a3b8;font-size:0.85rem}
footer p{margin-bottom:4px}
</style>
</head>
<body>
<div class="container">
<a href="index.html" class="back-link">← All Reports</a>
<header>
<h1>AEO &amp; SEO Daily Digest</h1>
<div class="subtitle">` + html.EscapeString(dateStr) + `</div>
</header>
`)

	if len(categories) == 0 {
		sb.WriteString(`<p style="color:#64748b;padding:32px 0">No SEO/AEO items found for this date.</p>`)
	}

	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf(`<details open>
<summary>%s <span class="item-count">%d items</span></summary>
<div class="items">
`, html.EscapeString(cat.Name), len(cat.Items)))

		for _, item := range cat.Items {
			scoreClass := "badge-score"
			if item.Score >= 80 {
				scoreClass = "badge-score-high"
			}

			sb.WriteString(fmt.Sprintf(`<div class="item">
<div class="item-title"><a href="%s" target="_blank" rel="noopener noreferrer">%s</a></div>
<div class="item-meta">
<span class="badge badge-source">%s</span>
`,
				html.EscapeString(item.URL),
				html.EscapeString(item.Title),
				html.EscapeString(item.Source),
			))

			for _, extra := range item.ExtraSources {
				sb.WriteString(fmt.Sprintf(`<span class="badge badge-source">%s</span>`, html.EscapeString(extra)))
			}

			sb.WriteString(fmt.Sprintf(`<span class="badge %s">Score: %d</span>`, scoreClass, item.Score))

			if !item.PublishedAt.IsZero() {
				sb.WriteString(fmt.Sprintf(`<span class="badge-time">%s</span>`, html.EscapeString(item.PublishedAt.Format("Jan 2, 15:04 UTC"))))
			}

			sb.WriteString("</div></div>\n")
		}

		sb.WriteString("</div></details>\n")
	}

	sb.WriteString(fmt.Sprintf(`<footer>
<p>Generated: %s</p>
<p>Sources: %s</p>
</footer>
</div>
</body>
</html>
`, html.EscapeString(generatedAt.UTC().Format(time.RFC3339)), html.EscapeString(strings.Join(sourceNames, ", "))))

	return sb.String()
}
