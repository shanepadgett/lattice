# Subplan 02: Config Load, Merge, Validate

## Goal

Implement config loading, deep merge, and validation to produce a canonical config.

## Inputs

- Base config JSON
- Optional site config JSON

## Steps

1. Define config structs and JSON parsing.
2. Implement deep merge with array replace semantics.
3. Validate required fields and referenced breakpoint names.
4. Normalize tokens into canonical maps (colors, scales, fonts).
5. Add a CLI command: lcss config print.

## Deliverables

- Merged config output (debug command).
- Deterministic ordering for serialized output.

## Acceptance checks

- Missing required fields produce explicit errors.
- Output shows merged values from base + site configs.
