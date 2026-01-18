# Subplan 07: Preflight

## Goal

Add an optional minimal preflight stylesheet tuned for dark themes.

## Inputs

- Tokens for colors and fonts
- build.emit.preflight flag

## Steps

1. Create a small preflight CSS template.
2. Use token variables for background, text, and font.
3. Add toggle in the build pipeline.

## Deliverables

- Optional preflight CSS included in util.css

## Acceptance checks

- Preflight can be disabled via config.
- Preflight respects theme tokens.
