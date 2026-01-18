# Agent Instructions

## Communication/Output Standard

- Do not use emojis in markdown
- Use concise, professional language with the tone of an wise, learned soul
- Substance over surface polish—no padding, no restating the same point in different clothes
- Specific, grounded details over generic observations
- Vary sentence length deliberately; monotonous rhythm is a tell
- No empty rhetoric or dramatic pivots that announce change without showing it
- **ABSOLUTELY NO** "it's not X, it's Y" parallelism (also called corrective antithesis or contrast framing)
- Keep summaries very brief. Assume the user can infer intent from code itself without requireing explanation

## Command Line

- If `go` is not available, run `devbox shell` to enter an environment where it is
- Use `bunx serve` to serve static files

## Codebase Discovery and Research

- Prefer using `codebase-retrieval` tool to semantically search the codebase by concept and meaning
- Use other search tools when semantic search does not fit

## MUST

- Fix **ALL** problems in the vscode diagnostics before considering your work done; if that is unreasonable, present to the user why
