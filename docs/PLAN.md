# Go Utility-CSS Engine (Internal Tailwind-like) — Project Plan

This plan describes a Go-based utility CSS generator driven by per-site JSON config. It targets internal use, minimal “framework” features, no plugin system, no light mode, and assumes all sites are dark-themed (with optional multiple dark themes).

---

## Goals

* **Single Go CLI** that builds CSS from:

  * a **base/default config** (shared across sites)
  * an optional **site override config** (colors, fonts, etc.)
  * the set of **utility classes actually used** in project files (JIT-style)
* **Stable class naming** across sites.
* **Simple variants** (responsive + state), explicitly **no `dark:` variant** (dark is default).
* **Multiple themes allowed**, but all themes are dark (theme is a token set, not a variant).

Non-goals (explicitly):

* No external plugin ecosystem
* No light mode support
* No arbitrary value syntax (`p-[13px]`) unless you decide later

---

## What you’ll build (high-level architecture)

### Inputs

1. **Base config JSON** (your “design system defaults”)
2. **Site config JSON** (optional override: colors, fonts, maybe spacing scale tweaks)
3. **Source files** (templates/components): `.html .tmpl .gohtml .tsx .jsx .vue .svelte .mdx` etc.

### Output

* `dist/lattice.css` (tokens as CSS variables + optional base + utilities)
* `dist/manifest.json` (optional debug: used classes, generation stats)

---

## Design choices to keep it simple and internal-friendly

### 1) Token-first: CSS variables everywhere

Generate tokens as CSS variables, then utilities reference them.

Example:

* `--space-4: 1rem;`
* `.p-4 { padding: var(--space-4); }`

This makes per-site overrides easy:

* site overrides only change variables (colors/fonts), utilities stay the same.

### 2) Themes are token sets (not variants)

You can support multiple dark themes by scoping token variables:

* Default theme:

  * `:root { --color-bg: ... }`
* Alternative theme:

  * `[data-theme="midnight"] { --color-bg: ... }`

Switching theme becomes runtime HTML attribute change, no class variant needed.

### 3) Minimal variant set

* Responsive: `sm:`, `md:`, `lg:`, `xl:` (config-driven)
* State: `hover:`, `focus:`, `active:`, `disabled:` (config-driven)
* Optional: `group-hover:` (later)

No `dark:`.

---

## Config format (JSON)

### Base config: `configs/base.json`

Suggested structure (you can adjust naming):

```json
{
  "schemaVersion": 1,
  "classPrefix": "",
  "separator": ":",
  "breakpoints": {
    "sm": "640px",
    "md": "768px",
    "lg": "1024px",
    "xl": "1280px"
  },
  "themes": {
    "default": {
      "colors": {
        "bg": "#0b0f14",
        "fg": "#e6edf3",
        "muted": "#9aa7b2",
        "primary": "#7aa2f7",
        "danger": "#f7768e"
      },
      "font": {
        "sans": "ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Arial",
        "mono": "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace"
      }
    }
  },
  "scales": {
    "space": {
      "0": "0",
      "1": "0.25rem",
      "2": "0.5rem",
      "3": "0.75rem",
      "4": "1rem",
      "6": "1.5rem",
      "8": "2rem",
      "12": "3rem",
      "16": "4rem"
    },
    "radius": {
      "0": "0",
      "2": "0.25rem",
      "4": "0.5rem",
      "8": "1rem"
    },
    "fontSize": {
      "sm": "0.875rem",
      "base": "1rem",
      "lg": "1.125rem",
      "xl": "1.25rem",
      "2xl": "1.5rem"
    },
    "lineHeight": {
      "tight": "1.1",
      "snug": "1.3",
      "normal": "1.5",
      "relaxed": "1.7"
    }
  },
  "variants": {
    "responsive": ["sm", "md", "lg", "xl"],
    "state": ["hover", "focus", "active", "disabled"]
  },
  "build": {
    "content": ["./src/**/*.{html,tmpl,gohtml,tsx,jsx,vue,svelte,mdx}"],
    "safelist": [],
    "unknownClassPolicy": "warn",
    "emit": {
      "tokensCss": true,
      "base": true,
      "manifest": true
    }
  }
}
```

### Per-site override: `configs/site.json`

Override only what you need:

```json
{
  "themes": {
    "default": {
      "colors": {
        "primary": "#a78bfa",
        "danger": "#fb7185"
      },
      "font": {
        "sans": "Inter, ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Arial"
      }
    },
    "midnight": {
      "colors": {
        "bg": "#070a10",
        "fg": "#dbe7ff",
        "primary": "#22d3ee"
      }
    }
  }
}
```

Merge behavior:

* Deep-merge objects
* Arrays: replace (or append if you explicitly want)

---

## Class grammar (keep it strict)

### Pattern

`[variant1:variant2:...]utility-token`

Examples:

* `p-4`
* `px-6`
* `text-base`
* `bg-primary`
* `md:p-8`
* `hover:bg-primary`

### Escaping rule

Keep it simple: only allow `a-zA-Z0-9-:_` in class names.
If you want to support more later, introduce an escaping scheme in v2.

---

## Utility set (initial)

Start small and expand. Suggested phase 1 utilities:

### Spacing

* padding: `p-*`, `px-*`, `py-*`, `pt-*`, `pr-*`, `pb-*`, `pl-*`
* margin: `m-*`, `mx-*`, `my-*`, `mt-*`, `mr-*`, `mb-*`, `ml-*`
* gap: `gap-*`, `gapx-*`, `gapy-*` (or `gap-x-*`/`gap-y-*`)

### Layout

* display: `block`, `inline-block`, `inline`, `flex`, `grid`, `hidden`
* width/height: `w-*`, `h-*` (scale + special: `full`, `screen` later)
* flex: `flex-row`, `flex-col`, `items-center`, `justify-between`, etc.

### Typography

* `font-sans`, `font-mono`
* `text-sm|base|lg|xl|2xl`
* `leading-tight|snug|normal|relaxed`
* `font-400|500|600|700` (if you want numeric weights)
* `tracking-*` (later)

### Colors

* `bg-*` for theme color keys: `bg-bg`, `bg-primary`, etc.
* `text-*` similarly
* `border-*`

### Borders / radius / shadow

* `border`, `border-0`, `border-2` (optional)
* `rounded-*` from radius scale
* shadows: `shadow-sm|md|lg` (if you define them)

### Interaction

* `cursor-pointer`, `select-none`
* `opacity-*` (scale)
* `transition`, `duration-*` (later)

---

## Implementation plan (milestones)

### Milestone 0: Repo skeleton

Create:

* `cmd/lcss/main.go` (CLI entry)
* `internal/config/` (load + merge + validate)
* `internal/extract/` (scan files for class candidates)
* `internal/compile/` (class -> CSS rules)
* `internal/emit/` (write output files)
* `internal/util/` (hashing, sorting, sets)

### Milestone 1: Config load + merge + validation

* Load base + site JSON
* Apply deep merge
* Validate required fields:

  * schemaVersion
  * themes.default exists
  * scales.space exists
  * breakpoints referenced by variants exist
* Normalize tokens into a canonical structure:

  * flatten `themes.<name>.colors.*` into `--color-<key>`
  * flatten scales into `--space-<key>`, `--radius-<key>`, etc.

Deliverable:

* `lcss config print --base configs/base.json --site configs/site.json` shows merged config.

### Milestone 2: Token CSS emission

Generate token CSS variables for inclusion in `lattice.css`:

* `:root` includes default theme vars
* For each additional theme:

  * `[data-theme="<name>"] { ...vars... }`

Deliverable:

* consistent output ordering (sorted keys)
* stable formatting

### Milestone 3: Class extractor (JIT input)

Implement a fast extractor:

* Glob content paths from config
* For each file, read text
* Extract candidate strings from:

  * `class="..."`
  * `className="..."`
  * `class='...'`
* Split on whitespace
* Add safelist entries
* Output a set of class strings

Deliverable:

* `lcss scan --config ...` prints count + top N classes

Notes:

* Support dynamic classes later via safelist patterns (regex/glob) if needed.

### Milestone 4: Core compiler for utilities

Define an internal registry of utility matchers.

A simple approach:

* Each utility “family” has a `Match(class string) (rule, ok)`
* `rule` contains:

  * selector (escaped)
  * declarations
  * optional media query wrapper
  * optional pseudo-class wrapper

Implement:

* spacing utilities (`p-`, `m-`, `gap-`)
* typography basics (`text-`, `font-`, `leading-`)
* colors (`bg-`, `text-`, `border-`)
* layout (`flex`, `grid`, `hidden`, alignment)
* radius (`rounded-`)

Deliverable:

* `lcss build` outputs `lattice.css` containing only rules for used classes.

### Milestone 5: Variants pipeline

Parse prefix chain:

* split by separator `:` (config-driven)
* last segment is base utility; earlier segments are variants

Apply:

* responsive variant -> wrap in `@media (min-width: X)`
* state variant -> append pseudo `:hover` etc.

Selector examples:

* `hover:bg-primary` -> `.hover\:bg-primary:hover { background-color: var(--color-primary); }`
* `md:p-8` -> `@media (min-width: 768px) { .md\:p-8 { padding: var(--space-8); } }`
* `md:hover:bg-primary` -> media wrapper + `:hover`

Deliverable:

* variant order handling: responsive wrapper should be outermost.

### Milestone 6: Base stylesheet (optional, minimal)

Since all sites are dark, base can be small and token-based:

* set `color-scheme: dark;`
* body background/foreground from tokens
* default font from tokens
* basic resets you like

Make it toggleable with `build.emit.base`.

### Milestone 7: Watch mode + caching

* Hash inputs (config files + file modification times + extracted class set)
* Only rebuild if changed
* `lcss watch` uses fsnotify

---

## CLI commands (suggested)

* `lcss build --base configs/base.json --site configs/site.json --out ./dist`
* `lcss watch --base ... --site ... --out ./dist`
* `lcss scan --base ... --site ...` (debug: show used classes)
* `lcss config print --base ... --site ...` (debug merged config)

---

## Distribution (internal)

### Recommended for internal sites you own

* Keep `configs/base.json` in a central repo (or a subdir you copy)
* Each site has its own `configs/site.json`
* Distribute the compiler as:

  * a GitHub Release binary (mac/win/linux)
  * optionally a Docker image for CI

If you want “one-liner installs”:

* provide a tiny script in your internal docs that downloads the right binary + verifies checksum.

---

## Testing strategy (lightweight but effective)

1. **Golden file tests** — Given (config + class set) -> `lattice.css` matches expected output.
1. **Extractor tests** — Ensure class extraction works for your templating stack.
1. **Config merge tests** — Ensure overrides behave (colors/fonts) without breaking base.

---

## Practical shortcuts (worth doing)

* **Deterministic output ordering**

  * Sort classes lexicographically before compiling
  * Sort declarations in a stable order
* **Rule de-dup**

  * If two classes would emit identical CSS (rare early), keep first.
* **Helpful errors**

  * Unknown token keys should warn (and optionally fail CI via `--strict`)

---

## Roadmap additions (only if you feel pain later)

* safelist patterns (globs/regex)
* `group-hover:` and `focus-within:`
* `@container` queries (if you want)
* arbitrary values (avoid unless truly needed)
* better template parsing for frameworks with heavy string interpolation

---

## Suggested “Day 1” build target

By end of first implementation day, aim for:

* `lcss build` that:

  * merges base+site config
  * emits `lattice.css`
  * scans `./src`
  * generates spacing + color + typography utilities for used classes
  * supports `md:` and `hover:`
* output is stable and checked into `dist/` for a sample site

---

## Repo layout example

```text
lcss/
  cmd/
    lcss/
      main.go
  internal/
    config/
      load.go
      merge.go
      validate.go
      normalize.go
    extract/
      extract.go
      patterns.go
    compile/
      compile.go
      variants.go
      escape.go
      utilities_spacing.go
      utilities_color.go
      utilities_typography.go
      utilities_layout.go
    emit/
      tokens.go
      utilcss.go
      manifest.go
    util/
      hash.go
      set.go
      sort.go
  configs/
    base.json
  README.md
  go.mod
```
