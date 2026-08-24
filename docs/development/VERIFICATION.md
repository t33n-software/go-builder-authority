# Go Builder Authority Verification

## Local source verification

The GBA-1 source foundation is verified without tenant credentials or builder
registry access:

```text
go run -mod=readonly ./cmd/build
go test -count=1 ./...
go run -mod=readonly ./cmd/check-coverage
```

The Go toolchain is pinned exactly (`toolchain go1.26.6`,
`GOTOOLCHAIN=local`); no lane downloads a toolchain at build time.

The source gate checks:

```text
Go formatting
go mod verify
go mod tidy -diff
tools module download, verify, and tidy -diff
staticcheck lint
unit tests
100% statement coverage
race detection
go vet
govulncheck fail-closed vulnerability analysis
Lefthook configuration validation
Linux/AMD64 source-gate build
embedded module provenance
```

Build tools (`govulncheck`, `staticcheck`, `lefthook`) live in the separate
pinned `tools/` module with its own verified `go.mod` and committed `go.sum`;
they never join the source module graph. The same module pins the canonical
quality gate and coverage gate of the go-quality-authority territory home and
the canonical conformance verifier of the repository-governance home. CI
re-runs the full gate on a daily schedule so newly disclosed vulnerabilities
in the pinned toolchain or dependency graph fail closed even without source
changes. Lefthook provides the local `commit-msg` hook (governed
commit-message validation) and the canonical pre-push validation through
`git-governance --interactive never validate pre-push`.

## CI gates

The shared-line workflows are the byte-identical canonical callers of the
repository-governance home, pinned by full-length commit SHA: `ci.yml` runs
the canonical quality gate of the go-quality-authority territory home (check
context `Quality gates / linux-amd64`), `codeql.yml` runs the canonical
CodeQL lane (check context `CodeQL / CodeQL (go)`, consumed by the
code-scanning rule-set rule), and `dependency-review.yml` runs the dependency
admission review (check context `Dependency review / Dependency admission
review`). The callers trigger on push and pull request to every shared line
(`main`, `develop`, `release/**`, `support/**`) plus a daily schedule and
manual dispatch. The `canonical-conformance.yml` workflow runs the home's
conformance verifier (check context `Canonical conformance`) against
`repo-bindings.json`: caller hashes and pins, canonical file equality,
CODEOWNERS materialization, config-seam conformance, tool-pin admission, and
license-lane wiring. The organization rule-sets bind a check context only
after the lane has proven it on a real pull request to the exact target line;
the binding of this repository is documented in
`docs/conventions/hosting-plattform/github/rule-sets/`.

## Deliberate delivery boundary

GBA-1 does not claim a released Go builder image. Builder production remains
blocked until a separate authority binds:

```text
approved internal Go proxy
approved base-image and toolchain inputs
builder artifact and evidence repositories
builder signer and promoter identities
SBOM, signature, provenance, scan, policy, approval, and revocation evidence
immutable builder promotion and consumer admission
```

No local cache, public registry fallback, tenant credential, or ad-hoc Docker
build may substitute these controls.
