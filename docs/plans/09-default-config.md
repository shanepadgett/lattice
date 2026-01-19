# Subplan 09: Default Config (ship-ready defaults)

## Goal

Define and ship a sensible default config JSON with a complete token set and scales, suitable for most sites out of the box.

## Inputs

- Utility compiler and tokens pipeline
- Current set of supported utilities and scales
- Desired defaults (Tailwind-like baseline, tuned for lattice)

## Steps

1. Inventory required scales based on supported utilities (space, size, radius, typography, effects, motion).
2. Draft base theme tokens (colors, fonts) for a dark-first default theme.
3. Define full scale maps with names and values for each supported scale.
4. Choose breakpoint defaults and variant defaults (responsive + state).
5. Produce a default config JSON with build settings and content globs.
6. Review defaults for consistency, naming, and completeness.

## Deliverables

- A ship-ready default config JSON at configs/default.json
- Documented rationale for key defaults (short notes)

## Acceptance checks

- Defaults cover every utility scale without missing tokens.
- Default config compiles and emits tokens + utilities without errors.
- Defaults are conservative and consistent for broad site usage.

## Implementation

- Default config: configs/default.json

## Default notes (concise)

- Dark-first theme tokens with a numeric palette (no semantic color names).
- Default font stacks use non-system fonts and rely on `fonts.imports` for loading.
- `radius.default` and `shadow.default` are present to support `rounded` and `shadow` utilities.
- `container.default` ensures `container` emits a sensible max width even without a size suffix.
- Size, max-width, and max-height scales are included to cover `w-*`, `max-w-*`, and `max-h-*` without relying solely on spacing.
- Motion scales keep modest defaults (short durations, common easings) for UI transitions.
