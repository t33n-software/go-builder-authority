# Go Builder Authority

`go-builder-authority` is the tenant-neutral authority for approved Go build
images. It produces no tenant application artifacts and holds no tenant
credentials, tenant runtime configuration, or tenant deployment authority.

## Authority boundary

The authority owns the lifecycle of a Go builder artifact:

```text
builder source and pinned inputs
-> builder artifact subject
-> SBOM, signature, provenance, policy, approval, and revocation evidence
-> immutable approved builder digest
-> controlled consumer authorization
```

Consumers may use a builder only through a complete immutable digest after
re-verifying its evidence. Runner, Docker, BuildKit, and Node caches are not
builder trust or availability authorities.

## Bootstrap status

The root commit establishes the repository boundary and directory layout.
GBA-1 adds source-quality, workflow, Ruleset, evidence, and builder-production
contracts through governed ticket work. It does not claim a released builder
artifact until the separately required infrastructure and evidence authorities
exist.

In CI the repository is a tenant of the canonical repo surface: the three
shared-line workflows (`ci.yml`, `codeql.yml`, `dependency-review.yml`) are
byte-identical callers of the repository-governance home, and the canonical
quality gate of the go-quality-authority territory home runs through the
tooling module. The `repo-bindings.json` manifest binds the adoption (home
pin, fleet classes, caller and file hashes, config-seam and tool-catalog
versions), and the `Canonical conformance` check re-proves it fail-closed on
every shared-line change.

## Repository layout

- `builder/go/` contains the Go builder definition and controlled input
  manifests.
- `policy/` contains versioned builder policy references without credentials.
- `.github/` contains the canonical shared-line workflow callers, the
  canonical conformance workflow, and the ownership contract.
- `cmd/` contains repository-local verification tooling.
- `internal/authority/` contains authority-owned Go logic when materialized.
- `internal/packaging/` contains whitebox workflow and packaging contracts.
- `repo-bindings.json` binds the canonical repo-surface adoption (home pin,
  fleet classes, caller and file hashes, config-seam and tool-catalog
  versions).
- `docs/` contains architecture, conventions, operations, specification, and
  development documentation.

## Non-goals

This repository must never contain:

- tenant App IDs, installations, private keys, allowlists, WIF providers, or
  runtime identities;
- tenant artifact, deployment, or release evidence;
- mutable builder tags, public builder fallbacks, or cache-based trust;
- a tenant credential broker runtime.
