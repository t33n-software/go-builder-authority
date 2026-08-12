# Go Builder Authority Verification

## Local source verification

The GBA-1 source foundation is verified without tenant credentials or builder
registry access:

```text
go run -mod=readonly ./cmd/build
go test -count=1 ./...
go run -mod=readonly ./cmd/check-coverage
```

The source gate checks:

```text
Go formatting
go mod verify
go mod tidy -diff
unit tests
100% statement coverage
race detection
go vet
Linux/AMD64 source-gate build
embedded module provenance
```

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
