package renderer

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// ReportMeta holds metadata about a report for the hub page.
type ReportMeta struct {
	Date           time.Time
	Filename       string // e.g. "2026-05-01.html"
	ItemCount      int
	CategoryCount  int
}

// RenderIndex renders the hub index.html listing all reports newest first.
func RenderIndex(reports []ReportMeta, generatedAt time.Time) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AEO &amp; SEO Daily Digest</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:'Inter',sans-serif;background:#f8fafc;color:#1e293b;line-height:1.6;padding:0 16px}
.container{max-width:900px;margin:0 auto;padding:32px 0}
header{margin-bottom:32px;padding-bottom:24px;border-bottom:2px solid #0d9488}
h1{font-size:1.75rem;font-weight:600;color:#0d9488}
.subtitle{color:#64748b;margin-top:8px;font-size:0.95rem}
.reports{list-style:none;margin-top:24px}
.report-item{background:#fff;border:1px solid #e2e8f0;border-radius:8px;margin-bottom:12px;padding:16px 20px;display:flex;align-items:center;gap:16px}
.report-item:hover{border-color:#0d9488}
.report-date a{font-weight:600;font-size:1rem;color:#0d9488;text-decoration:none}
.report-date a:hover{text-decoration:underline}
.report-stats{margin-left:auto;display:flex;gap:8px}
.stat{background:#f1f5f9;color:#64748b;font-size:0.8rem;font-weight:500;padding:3px 10px;border-radius:12px}
.empty{color:#94a3b8;padding:32px 0;text-align:center}
footer{margin-top:48px;padding-top:24px;border-top:1px solid #e2e8f0;color:#94a3b8;font-size:0.85rem}
</style>
</head>
<body>
<div class="container">
<header>
<h1>AEO &amp; SEO Daily Digest</h1>
<div class="subtitle">Daily aggregation of AEO and SEO news, categorized and ranked by impact</div>
</header>
`)

	if len(reports) == 0 {
		sb.WriteString(`<div class="empty">No reports yet. Run <code>seo-report generate</code> to create the first report.</div>`)
	} else {
		sb.WriteString(`<ul class="reports">`)
		for _, r := range reports {
			sb.WriteString(fmt.Sprintf(`<li class="report-item">
<div class="report-date"><a href="%s" target="_blank" rel="noopener noreferrer">%s</a></div>
<div class="report-stats">
<span class="stat">%d items</span>
<span class="stat">%d categories</span>
</div>
</li>
`,
				html.EscapeString(r.Filename),
				html.EscapeString(r.Date.Format("January 2, 2006")),
				r.ItemCount,
				r.CategoryCount,
			))
		}
		sb.WriteString(`</ul>`)
	}

	sb.WriteString(fmt.Sprintf(`<footer>
<p>Last updated: %s</p>
</footer>
</div>
</body>
</html>
`, html.EscapeString(generatedAt.UTC().Format(time.RFC3339))))

	return sb.String()
}
