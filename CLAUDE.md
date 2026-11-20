# CLAUDE.md

All project rules live in **[AGENTS.md](AGENTS.md)**. Read it before any change, then read [tasks.md](tasks.md) to find the current task.

Claude-specific notes:

- Before starting a task, judge whether the session's remaining usage can finish the whole slice. If it cannot, split the task in tasks.md or stop, rather than leave work half done.
- For UI work, use the `/baseline-ui` skill. Before each commit of code, run `/simplify` on the diff.
- Requirements come from the Claude Doc linked at the top of [docs/requirements.md](docs/requirements.md). If the doc has changed, re-export it instead of editing the snapshot by hand.
- Commit messages carry no `Co-Authored-By` or other AI attribution trailers.
