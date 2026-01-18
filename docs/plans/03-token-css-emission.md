# Subplan 03: Token CSS Emission

## Goal

Generate tokens.css from the canonical config with stable ordering.

## Inputs

- Canonical config from Subplan 02

## Steps

1. Convert tokens to CSS variables in :root.
2. Emit additional theme scopes as [data-theme="theme-name"] blocks.
3. Ensure sorted output ordering for diffs and caching.
4. Add CLI command: ucss tokens (or as part of build).

## Deliverables

- dist/tokens.css
- Optional: merged tokens into util.css when configured

## Acceptance checks

- Default theme variables appear in :root.
- Additional theme variables scoped correctly.
