---
agent: agent
---

# Conventional Commits

- Analyze the changes and determine their semantic intent
- Produce a single, concise Conventional Commit message (one sentence) that summarizes the change
- Do not mention specific files in commit messages
- Group edits into coherent, self-contained concerns (things that must change together to work), but avoid over-splitting into many tiny file-by-file commits
- Use multiple commits only for distinct, unrelated concerns
- Prefer `git add -A && git commit -m "<message>"` for one commit, or add specific files when committing separately
- Write messages that describe what the change enables or fixes for changelogs

## Execute Commit
- You **MUST** execute the commit(s) without explaining your plan to the user
