---
name: release
description: Use when releasing a daily seo-report digest — runs the full pipeline from fetching sources through generating the report, committing, opening a PR, merging, and creating a versioned GitHub release.
---

# release

Runs the full seo-report release pipeline: fetch → classify → render → commit → PR → merge → GitHub release.

## One-command release

```bash
make release
```

This chains: `generate` → `commit` → `pr` → `merge` → `gh-release`

## Step-by-step (if any step needs re-running)

| Step | Command | What it does |
|---|---|---|
| 1. Build | `make build` | Compile `./seo-report` binary |
| 2. Generate | `make generate` | Fetch RSS, classify, score, dedup, write `docs/` |
| 3. Commit | `make commit` | Create `chore/seo-report-YYYY-MM-DD` branch, commit `docs/` |
| 4. PR | `make pr` | Push branch, open draft PR |
| 5. Merge | `make merge` | Mark PR ready, squash-merge, delete branch, pull `main` |
| 6. Release | `make gh-release` | Tag `vYYYY-MM-DD`, create GitHub release with report link |

## Guards

- `make commit` fails if `docs/` has no changes — run `make generate` first
- `make pr` fails if already on `main` — run `make commit` first
- `make merge` fails if no open PR found on current branch

## Preview without releasing

```bash
make dry-run   # prints HTML to stdout, writes nothing
```

## Common mistakes

**Nothing to commit** — `docs/` already contains today's report. If re-running on same day, delete `docs/YYYY-MM-DD.html` first and re-run `make generate`.

**PR already merged** — `make merge` will error if PR is closed. Check `gh pr list` and proceed from the appropriate step.

**Wrong date** — `DATE` is evaluated at `make` invocation time using `date -u`. All targets in the same `make release` run use the same date.
