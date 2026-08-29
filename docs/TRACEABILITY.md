# Traceability

## GBA-1: Establish Go Builder Authority

Status: in implementation.

GBA-1 establishes:

```text
tenant-neutral Go Builder Authority boundary
main and develop source foundation
official GBA-1 feature branch
Go source-quality and exact-coverage gates
CodeQL, dependency admission, Dependabot, and Lefthook contracts
repository-specific working-branch, develop, and main Ruleset sources
builder authority, evidence, policy, and consumer-boundary documentation
```

GBA-1 does not establish:

```text
an approved Go builder image
an internal Go proxy
builder artifact or evidence registries
builder signer, promoter, or revocation identities
builder production, promotion, release, or support-line workflows
tenant configuration, credentials, or deployment infrastructure
```

These omissions are deliberate fail-closed boundaries. A builder artifact
cannot become approved or be consumed until each missing authority and evidence
contract is explicitly implemented and verified.

## GBA-3: Migrate to the t33n-software organization namespace

Status: in implementation.

GBA-3 establishes:

```text
module path and GitHub setup references on the t33n-software organization
  namespace
push-protections Ruleset source 00-push-protections.json in the verified
  GitHub export format
```

GBA-3 does not establish:

```text
an approved Go builder image
builder artifact or evidence registries
any change to the builder authority, evidence, policy, or consumer contracts
any change to the GBA-2 builder artifact delivery branch
```

## GBA-4: Align the Go 1.26.6 toolchain and source gates

Status: in implementation.

GBA-4 establishes:

```text
Go toolchain pin go1.26.6 with GOTOOLCHAIN=local
pinned tools module with govulncheck, staticcheck, and Lefthook
fail-closed vulnerability analysis in the source gate
Lefthook configuration validation and commit-msg hook
daily CI re-scan of the full source gate
```

GBA-4 does not establish:

```text
an approved Go builder image
builder artifact or evidence registries
any change to the builder authority, evidence, policy, or consumer contracts
any change to the GBA-2 builder artifact delivery branch
```

## GBA-6: Adopt the canonical repo surface

Status: in implementation.

GBA-6 establishes:

```text
schema-v3 quality configuration with the controlled toolchain identity
version surfaces on the development tools
canonical tool pins (go-quality-authority v1.0.1, repository-governance
  verifier)
canonical file family and the materialized CODEOWNERS contract
byte-identical canonical workflow callers and the tenant binding manifest
canonical conformance check
```

GBA-6 does not establish:

```text
an approved Go builder image
builder artifact or evidence registries
any change to the builder authority, evidence, policy, or consumer contracts
any change to the GBA-2 builder artifact delivery branch
```

## GBA-8: Reference the canonical gate chain through the tool pin

Status: in implementation.

GBA-8 establishes:

```text
the canonical quality-gate invocation through the pinned tooling module
the schema-owned default ticket-family scope
the removal of repository-local gate-chain copies
the fail-closed contract proof that no local gate-chain directory remains
```

GBA-8 does not establish:

```text
an approved Go builder image
builder artifact or evidence registries
any change to the builder authority, evidence, policy, or consumer contracts
any change to the GBA-2 builder artifact delivery branch
```
