# Subplan 04: Class Extractor

## Goal

Scan content files to extract used class names for JIT compilation.

## Inputs

- Config content globs
- Safelist entries

## Steps

1. Resolve globs and read file contents.
2. Extract candidate class strings from class/className attributes.
3. Split by whitespace and add safelist entries.
4. Return a unique, sorted set of class names.
5. Add CLI command: lcss scan.

## Deliverables

- In-memory class set for the build pipeline.
- Optional debug output with counts.

## Acceptance checks

- Correctly extracts classes from HTML-like attributes.
- Ignores invalid characters per class grammar.
