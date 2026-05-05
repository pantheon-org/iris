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

## Provider testdata fixtures

Every provider must have two fixture files in `internal/providers/testdata/`:

- `<provider>_input.{json,toml}` — realistic on-disk config as a user would have it (source of truth for `Parse` tests).
- `<provider>_expected.{json,toml}` — exact output iris produces when it round-trips that input through `Parse` then `Generate` (source of truth for `Generate` tests).

The canonical test pattern:

```go
func TestXxxProvider_Parse_ExtractsServersFromFixture(t *testing.T) {
    content, _ := os.ReadFile("testdata/xxx_input.json")
    parsed, err := p.Parse(string(content))
    // assert specific fields on parsed servers
}

func TestXxxProvider_Generate_FixtureMatch(t *testing.T) {
    content, _ := os.ReadFile("testdata/xxx_input.json")
    servers, _ := p.Parse(string(content))
    got, _ := p.Generate(servers, string(content))
    expected, _ := os.ReadFile("testdata/xxx_expected.json")
    assert got == string(expected)
}
```

To regenerate expected files after a deliberate format change, run:

```sh
go test -tags gen_expected ./internal/providers/... -run TestGenExpected -v
```

(The generator lives in `internal/providers/gen_expected_test.go` — create it locally, do not commit it.)

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
