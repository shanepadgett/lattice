# Agent Instructions (Product)

Keep this file short and practical. Document only what an agent needs to use the tool correctly.

## Config

- Merge and validate config:
  - `lcss config print --base <path> [--site <path>]`
  - Writes canonical config JSON to stdout.

## Tokens

- Emit tokens CSS:
  - `lcss tokens --base <path> [--site <path>] [--out <path>]`
  - Use `--stdout` to write to stdout.

## Build

- Build lattice.css from config + content scan:
  - `lcss build --base <path> [--site <path>] [--out <path>] [--stdout]`
  - Uses `build.content` and `build.safelist` from config.

## Class Scan

- Extract used classes:
  - `lcss scan --base <path> [--site <path>] [--top <n>] [--per-file]`
- What you get:
  - Total files scanned and unique class count.
  - Top N classes by frequency.
  - Optional per-file top N when `--per-file` is set.

## Extraction Rules (Strict)

- Finds `class` and `className` values in HTML-like attributes.
- In Go templates, also scans `class={{ ... }}` and `className={{ ... }}` for string literals.
- Class names must match: `a-zA-Z0-9-:_`.
