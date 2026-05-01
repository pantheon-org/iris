# Contributing to iris

## Development setup

Install [mise](https://mise.jdx.dev/):

```sh
curl https://mise.run | sh
```

Then install tools and verify:

```sh
mise install
mise run build
mise run test
mise run lint
```

Install git hooks:

```sh
lefthook install
```

## Workflow

- Always work on a feature branch — never commit directly to `main`.
- Write tests before implementation (TDD).
- Run `go mod tidy` after adding/removing dependencies.
- Use `errors.Is` / `errors.As` for error checks — never string-match errors.
- All logic lives in `internal/`; `cmd/iris/main.go` only wires cobra commands.

## Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` new feature
- `fix:` bug fix
- `chore:` tooling / dependencies / CI
- `docs:` documentation only
- `test:` tests only
- `refactor:` no behaviour change

## Opening a PR

Fill in the PR template. CI must be green before merge.
Squash-merge is enforced; branches are deleted after merge.
