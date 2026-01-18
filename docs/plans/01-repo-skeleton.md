# Subplan 01: Repo Skeleton

## Goal

Create the initial Go CLI layout and internal package structure.

## Inputs

- Main plan Milestone 0

## Steps

1. Create cmd/ucss/main.go and wire a basic CLI entry.
2. Create internal packages: config, extract, compile, emit, util.
3. Add placeholder types/functions to allow compiling.
4. Add a README note on how to run the CLI (temporary).

## Deliverables

- CLI builds and runs with a placeholder command.
- Folder structure matches the plan.

## Acceptance checks

- go test ./... passes (even if empty tests).
- Running the CLI prints a help message.
