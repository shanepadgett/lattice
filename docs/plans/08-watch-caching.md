# Subplan 08: Watch Mode and Caching

## Goal

Implement watch mode and incremental rebuilds based on input changes.

## Inputs

- Config files
- Content files
- Extracted class set

## Steps

1. Hash config inputs and file modification state.
2. Skip rebuild if inputs are unchanged.
3. Implement watch mode using fsnotify.
4. Trigger rebuild on file or config changes.

## Deliverables

- ucss watch command with fast incremental builds

## Acceptance checks

- No rebuild when nothing changes.
- Rebuild triggers on config or content changes.
