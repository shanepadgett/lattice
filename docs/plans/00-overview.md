# Subplan 00: Overview

## Purpose

Break the project into incremental, testable milestones with clear deliverables and acceptance checks.

## Key constraints from the main plan

- Single Go CLI.
- Token-first (CSS variables) and stable class naming.
- Only responsive + state variants (no `dark:`).
- Theme switching via scoped token sets.
- Strict class grammar and JIT extraction.

## Dependencies

- Base + site config JSON files (example shapes in the main plan).
- Target file extensions for class extraction.

## Output artifacts

- dist/util.css
- dist/tokens.css
- dist/manifest.json (optional)

## Acceptance criteria

- Each milestone produces a working CLI command or artifact.
- Output order is stable and deterministic.
- No unused classes are emitted unless safelisted.
