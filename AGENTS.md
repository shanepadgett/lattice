# Agent Instructions

## Communication/Output Standard

- Do not use emojis in markdown
- Use concise, professional language with the tone of an wise, learned soul
- Substance over surface polish—no padding, no restating the same point in different clothes
- Specific, grounded details over generic observations
- Vary sentence length deliberately; monotonous rhythm is a tell
- No empty rhetoric or dramatic pivots that announce change without showing it
- **ABSOLUTELY NO** "it's not X, it's Y" parallelism (also called corrective antithesis or contrast framing)
- Keep summaries very brief when requested; do not label them as summaries
- Never describe how the answer will be structured; just give the answer
- When asking questions, always include 2–3 options and a clear recommendation

## Planning

- When planning, be thorough. Do not leave any decisions open, and present the user with actual implemention thoughts and details. Ask questions of the user if any part of the task at hand is ambiguous or has multiple paths

## Command Line

- If `go` is not available, run `devbox shell` to enter an environment where it is
- `gnumake` (GNU Make) is available in devbox; prefer `make build` for local builds
- Use `bunx serve` to serve static files
- Build (preferred): `make build`  # outputs to `./bin/lcss`
- Build (fallback): `go build -o ./bin/lcss ./cmd/lcss`
- Use `gofmt -w .` before providing summary if you edited Go files.

## MUST

- Fix **ALL** problems in the vscode diagnostics before considering your work done; if that is unreasonable, present to the user why
- **ONLY** update STATE.md when user has validated work is complete for a plan
- Validate all work with formatting and making sure it builds
- Whenever you implement something that exposes a new interface, cli command, build tool, anything, you must test out the functionality and fix any issues before considering the work complete
- When adding agent-relevant features, update docs/AGENT-INSTRUCTIONS.md with concise usage guidance
- Always clean up after any test artifacts are created unless the user is required to verify them

## MUST NOT

- Add tests. There will be no tests in this repository
- Introduce any new libraries. This repo is library free and we use standard library only.
