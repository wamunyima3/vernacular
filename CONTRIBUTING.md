# Contributing

Thanks for helping build open, verified data for Zambian languages. There are two ways to contribute.

## Contributing language data

Words, phrases, examples and reviews are contributed through the web app's contributor and reviewer tools, not through pull requests. Reviewers are native speakers approved per language (see DQ-3 in [docs/requirements.md](docs/requirements.md)). Do not open PRs that add dictionary data. Data enters only through a registered, licensed source or through the review pipeline.

## Contributing code

1. Read [AGENTS.md](AGENTS.md). It holds the conventions, the definition of done and the commit rules.
2. Pick an unchecked task in [tasks.md](tasks.md), or open an issue for work that is not listed yet.
3. Set up locally (see the README quickstart), then branch from `main`: `feat/<task-id>-short-name`.
4. Build the task as a vertical slice: contract → db → api → web → verify.
5. Run `make check` before pushing. Git hooks (lefthook) run formatters and linters on staged files, and check the commit message format.
6. Open a PR against `main` using the checklist below.

### PR checklist

- [ ] Links the task ID and the user story or requirement IDs
- [ ] Contract changed first; `make gen` run; generated code committed
- [ ] Tests for new behaviour, including failure cases
- [ ] `make check` passes locally
- [ ] Screenshots for UI changes, at phone and desktop widths
- [ ] Any new dependency names its licence and why an existing one did not fit
- [ ] tasks.md updated

## Licence

Code is MIT-licensed ([LICENSE](LICENSE)). The dataset will have its own open licence (CC BY 4.0 proposed; see §12 of the requirements). By contributing, you agree that your contribution is released under these licences.
