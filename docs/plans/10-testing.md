# Subplan 10: Testing (unit + integration)

## Goal

Add a test strategy and concrete tests (unit, integration, and CI) to validate behavior end-to-end using fixtures and deterministic outputs.

## Inputs

- Package code (config, extract, compile, emit)
- Example configs (`configs/base.json`, `configs/site.json`)
- Fixture HTML files representing common and edge case class usage

## Strategy

- Unit tests: cover parsing, merging, normalization, token generation, and individual matcher logic.
- Integration tests (fixture-driven): run the CLI or internal pipeline against small fixture projects and assert generated `lattice.css` and `manifest.json` contents match golden files.
- Snapshot/golden files: store expected outputs under `internal/*/testdata/golden` or a top-level `testdata/golden/<fixture-name>/` for easy diffs.
- Smoke tests: CLI `lcss scan`/`lcss build` run in a temporary directory to ensure commands exit successfully and produce files.
- Determinism checks: assert stable ordering and formatting (use canonical serializers) so golden diffs are meaningful.

## Steps

1. Create `testdata/` fixtures for the following cases:
   - Minimal HTML (single class)
   - Multiple classes with variants (e.g., `md:hover:bg-primary p-4`)
   - Safelist and invalid class names
   - Theme overrides (site config + fixture)
2. Implement unit tests for:
   - Config load + merge + validation
   - Token CSS generator (check variable names and scopes)
   - Class extractor (ensuring valid classes are captured)
   - Utility matchers (spacing, colors, variants)
3. Implement integration tests that:
   - Copy a fixture into a temp dir
   - Run the `lcss` build (or call internal API directly)
   - Compare generated outputs to golden files (textual diffs)
4. Add tests for edge cases:
   - Unknown classes (ignored or reported)
   - Variant ordering (responsive outermost)
   - Watch mode trigger (simulate file change if feasible)
5. Add CI workflow step to run `go test ./...` and any additional test commands.

## Deliverables

- `internal/*` tests exercising core logic
- `testdata/fixtures` and `testdata/golden` directories with example cases
- CI workflow entries running tests

## Acceptance checks

- `go test ./...` passes in CI
- Integration tests validate tokens, utilities, and variants via golden diffs
- Tests are deterministic and fast (aim for <1s/unit, <2s/integration where possible)

## Example test case (integration)

- Fixture: `fixture-basic/` with `index.html` containing `<div class="p-4 bg-primary md:hover:bg-primary">` and a site config overriding `primary` color
- Test: run build and assert `lattice.css` contains `--color-primary` with overridden value and `.p-4` with `padding: var(--space-4)` and a `@media (min-width: ...)` wrapper for `md:hover:bg-primary`.

---

Tips:

- Use `t.Run` subtests and `cmp.Diff` or similar helpers for readable failures.
- Keep golden files small and human-readable; prefer CSS formatted with one declaration per line for easier diffs.
- For watch-mode tests, manipulate files in a temp dir and wait for the build to produce expected output (with a short timeout).
