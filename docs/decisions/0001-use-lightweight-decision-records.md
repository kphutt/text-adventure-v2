# Use lightweight decision records

## Status
Accepted

## Context
Need a way to capture why decisions were made so future-me (or teammates) don't re-debate settled questions. Must be low friction — append-only, no maintenance burden.

## Decision
Use numbered markdown files in `docs/decisions/`. One decision per file, minimal template (status, context, decision, consequences). Never edit after writing — supersede with a new record if a decision changes.

## Consequences
Every repo gets the same `docs/decisions/` folder. Easy to grep across projects. No tooling required.
