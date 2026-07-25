# Task Completion

Before considering a coding task done, run `mise run pre-commit` (format + lint + test).
This is also what the user's global CLAUDE.md instructs ("Run tidy and test tasks before committing, when available").

If only test files changed, commit type should be `test`; if only GitHub workflow/action files changed, commit type should be `ci` (per user's global preferences, not repo-specific but applies here).

`golang:format`/`golang:lint` both depend on `golang:mod-tidy`/`mod-tidy-check` — so go.mod/go.sum drift will fail CI (`mise run ci` = lint + test) even if code compiles.
