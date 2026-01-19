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
- Base stylesheet is enabled by default; disable with `build.emit.base: false`.

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
- Class names must match: `a-zA-Z0-9-:_/%`.

## Utilities (Site-Ready)

- Layout: `block`, `inline`, `flex`, `grid`, `hidden`, `contents`.
- Positioning: `relative`, `absolute`, `fixed`, `sticky`, `inset-*`, `top-*`, `right-*`, `bottom-*`, `left-*`.
- Sizing: `w-*`, `h-*`, `min-w-*`, `min-h-*`, `max-w-*`, `max-h-*`, `container`.
- Spacing: `p*`, `m*`, `gap-*`, `gap-x-*`, `gap-y-*`.
- Flex/grid: `flex-*`, `items-*`, `justify-*`, `content-*`, `self-*`, `grid-cols-*`, `grid-rows-*`, `col-span-*`, `row-span-*`.
- Typography: `text-*` (size/color/align), `font-*`, `leading-*`, `tracking-*`, `uppercase`, `underline`, `italic`, `list-*`.
- Color: `bg-*`, `text-*`, `border-*`.
- Borders: `border`, `border-*` (width/style), `border-x-*`, `border-y-*`, `border-t-*`, `border-r-*`, `border-b-*`, `border-l-*`.
- Radius: `rounded`, `rounded-*`, `rounded-t|b|l|r|tl|tr|bl|br`.
- Effects: `shadow*`, `opacity-*`.
- Overflow/visibility: `overflow-*`, `visible`, `invisible`, `sr-only`.
- Object/aspect: `object-*`, `aspect-*`.
- Transitions: `transition*`, `duration-*`, `ease-*`, `delay-*`.
- Transforms: `translate-x-*`, `translate-y-*`, `rotate-*`, `scale-*`.
- Interaction: `cursor-*`, `pointer-events-*`, `select-*`, `isolate`.

## Token Scales

- Core scales: `space`, `size`, `radius`, `borderWidth`, `fontSize`, `lineHeight`, `fontWeight`, `letterSpacing`.
- Effects: `shadow`, `opacity`.
- Layout: `z`, `aspect`, `maxWidth`, `maxHeight`, `container`.
- Motion: `duration`, `easing`, `delay`, `translate`, `rotate`, `scale`.
