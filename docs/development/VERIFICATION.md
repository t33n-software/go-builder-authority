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
they never join the source module graph. CI re-runs the full gate on a daily
schedule so newly disclosed vulnerabilities in the pinned toolchain or
dependency graph fail closed even without source changes. Lefthook provides
the local `commit-msg` hook (governed commit-message validation) and the
pre-push source-quality gate.

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
