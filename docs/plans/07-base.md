# Subplan 07: Base stylesheet

## Goal

Add an optional minimal base stylesheet tuned for dark themes.

## Inputs

- Tokens for colors and fonts
- build.emit.base flag

## Steps

1. Create a small base CSS template.
2. Use token variables for background, text, and font.
3. Add toggle in the build pipeline.

## Deliverables

- Optional base CSS included in lattice.css

## Acceptance checks

- Base styles can be disabled via config.
- Base styles respect theme tokens.
