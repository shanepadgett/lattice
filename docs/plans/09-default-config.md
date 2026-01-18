# Subplan 09: Default Config (ship-ready defaults)

## Goal

Define and ship a sensible default base config JSON with a complete token set and scales, suitable for most sites out of the box.

## Inputs

- Utility compiler and tokens pipeline
- Current set of supported utilities and scales
- Desired defaults (Tailwind-like baseline, tuned for lattice)

## Steps

1. Inventory required scales based on supported utilities (space, size, radius, typography, effects, motion).
2. Draft base theme tokens (colors, fonts) for a dark-first default theme.
3. Define full scale maps with names and values for each supported scale.
4. Choose breakpoint defaults and variant defaults (responsive + state).
5. Produce a base config JSON with build settings and content globs.
6. Review defaults for consistency, naming, and completeness.

## Deliverables

- A ship-ready default base config JSON (location TBD, likely configs/base.json)
- Documented rationale for key defaults (short notes)

## Acceptance checks

- Defaults cover every utility scale without missing tokens.
- Base config compiles and emits tokens + utilities without errors.
- Defaults are conservative and consistent for broad site usage.
