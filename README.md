# seo-report

Daily AEO/SEO news digest as a GitHub Pages static site. Fetches RSS feeds from major SEO publications, classifies items by category, scores them by impact, deduplicates, and renders a styled HTML report.

## How it works

```
RSS feeds → classify → score → dedup → docs/YYYY-MM-DD.html + docs/index.html
```

GitHub Actions runs this daily at 08:00 UTC and commits the output to `docs/`.

## Setup

### 1. Clone and build

```bash
git clone https://github.com/YOUR_USERNAME/seo-report.git
cd seo-report
go build -o seo-report ./cmd/seo-report
```

### 2. Enable GitHub Pages

In your repo settings → **Pages** → Source: `docs/` folder on `main` branch.

### 3. That's it

The workflow runs automatically. Reports appear at `https://YOUR_USERNAME.github.io/seo-report/`.

## CLI usage

```bash
# Generate today's report and regenerate the index
./seo-report generate

# Generate report for a specific date
./seo-report generate --date 2026-05-01

# Dry run — print HTML to stdout without writing files
./seo-report generate --dry-run

# List configured RSS sources
./seo-report sources
```

## Configuration (`config.yaml`)

### Sources

```yaml
sources:
  rss:
    - name: Search Engine Journal
      url: https://www.searchenginejournal.com/feed/
      weight: 3   # 1–3; affects impact score
```

Weight affects the impact score: weight 2 adds +10, weight 3 adds +20.

### Categories

Defined as an ordered list with keywords and priority (1 = highest):

```yaml
categories:
  - name: "AI Overviews & AEO"
    keywords: ["ai overview", "answer engine", "aeo"]
    priority: 1
  - name: "Industry News"
    keywords: []   # catch-all
    priority: 9
```

Items are matched to the highest-priority category whose keywords appear in the title or description. Items matching no specific category fall through to "Industry News".

### SEO relevance filter

```yaml
keywords:
  seo_relevance:
    - SEO
    - AEO
    - search ranking
    # ...
```

Items must match at least one relevance keyword to be included at all.

## Impact scoring (1–100)

| Signal | Points |
|---|---|
| Base | 50 |
| Source weight 2 | +10 |
| Source weight 3 | +20 |
| Same story covered by another source | +5 (max +20) |
| Urgency keyword in title ("breaking", "just announced", "critical", "major update") | +10 |
| AI Overviews & AEO category | +10 |

Items within each category are sorted by score descending.

## Output

- `docs/YYYY-MM-DD.html` — per-day report, self-contained HTML
- `docs/index.html` — hub listing all reports, newest first, regenerated on every run

## Development

```bash
go test ./...          # run all tests
go build ./...         # verify compilation
```

## Adding a new source

Add an entry to `config.yaml` under `sources.rss` with a `name`, `url`, and `weight` (1–3). No code changes needed.
