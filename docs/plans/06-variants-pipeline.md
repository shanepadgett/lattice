# Subplan 06: Variants Pipeline

## Goal

Apply responsive and state variants to base utility rules.

## Inputs

- Parsed class parts (variant chain + utility)
- Breakpoints and variant lists from config

## Steps

1. Split class by separator and validate variants.
2. Apply state variants as pseudo-classes.
3. Wrap responsive variants in media queries.
4. Enforce outermost media query wrapping.

## Deliverables

- Variant-aware CSS rule generation

## Acceptance checks

- hover:bg-primary produces a hover selector.
- md:p-8 wraps in a min-width media query.
- md:hover:bg-primary nests in media query with :hover selector.
