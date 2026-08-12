# Contributing

## Branch and commit conventions

GBA-1 and later work use the canonical governed branch families:

```text
feature/<ticket>-<slug>
fix/<ticket>-<slug>
docs/<ticket>-<slug>
refactor/<ticket>-<slug>
chore/<ticket>-<slug>
test/<ticket>-<slug>
perf/<ticket>-<slug>
hotfix/<ticket>-<slug>
```

Commits use Conventional Commits with the bound ticket scope, for example:

```text
feat(GBA-1): establish Go Builder Authority source gates
```

`main` and `develop` are shared lines. They are never direct development
targets after the one-time root bootstrap. Published official working branches
are append-only and must not be routinely rebased, amended, or force-pushed.

## Required local verification

For every Go change:

```text
gofmt
go test ./...
go run -mod=readonly ./cmd/check-coverage
git diff --check
```

The source-quality gate additionally runs race detection, static analysis, and
the Linux/AMD64 source-gate build.

## Scope boundary

Do not add tenant credentials, tenant runtime settings, private keys, tokens,
mutable builder tags, or public-network fallback paths to this repository.
