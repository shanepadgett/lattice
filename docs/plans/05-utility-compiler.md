# Subplan 05: Utility Compiler

## Goal

Map class names to CSS rules and emit util.css.

## Inputs

- Class set from extractor
- Canonical config tokens and scales

## Steps

1. Implement utility matcher registry.
2. Add matchers for spacing, typography, colors, layout, and radius.
3. Emit rules only for matched classes.
4. Ensure selector escaping for variant prefixes.

## Deliverables

- dist/util.css with only used utilities

## Acceptance checks

- Unknown classes are ignored or reported (configurable).
- Output is deterministic and minimal.
