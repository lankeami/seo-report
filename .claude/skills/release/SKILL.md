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

## Re-running on the same date

Before retrying, assess current state and resume from the right step:

```bash
DATE=$(date -u +%Y-%m-%d)
git branch --show-current
gh pr list --head chore/seo-report-$DATE
gh release view v$DATE 2>/dev/null && echo "release exists"
```

| State | Action |
|---|---|
| On `main`, no branch yet | `make release` (full run) |
| On `main`, branch exists, no PR | `git checkout chore/seo-report-$DATE && make pr && make merge && make gh-release` |
| On release branch, no PR | `make pr && make merge && make gh-release` |
| PR open, not merged | `make merge && make gh-release` |
| PR merged, no release | `make gh-release` |
| Release tag already exists | `gh release delete v$DATE --yes && git tag -d v$DATE && git push origin --delete refs/tags/v$DATE && make gh-release` |
| Need to regenerate the report | `git checkout main && git branch -D chore/seo-report-$DATE && git push origin --delete chore/seo-report-$DATE && rm docs/$DATE.html && make release` |

## Common mistakes

**Nothing to commit** — `docs/` already contains today's report. See "Re-running on the same date" above.

**PR already merged** — `make merge` will error if PR is closed. Check `gh pr list` and proceed from the appropriate step.

**Wrong date** — `DATE` is evaluated at `make` invocation time using `date -u`. All targets in the same `make release` run use the same date.
