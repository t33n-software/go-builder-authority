# ADR-0001: Go Builder Authority

## Status

Accepted.

## Context

A Go build image is a supply-chain artifact, not a Dockerfile `FROM` value or
an ephemeral cache entry. The credential-broker platform and every other Go
consumer require an immutable, independently verified builder digest before an
isolated build can begin.

Producing that builder inside an application, tenant, release, or hotfix lane
would mix builder trust with consumer authority and permit an unreviewed
builder change to alter an unrelated artifact build.

## Decision

This repository is the bounded context and supply-chain authority for approved
Go builder artifacts.

It owns:

- Go builder source definitions and pinned build inputs;
- builder artifact, evidence, policy, approval, and revocation contracts;
- immutable builder-digest publication and controlled consumer admission.

It does not own:

- tenant configuration, secrets, identities, runtimes, or deployments;
- application or credential-broker platform release artifacts;
- consumer promotion or deployment decisions.

The resulting lifecycle is:

```text
builder source, base image, and toolchain
-> controlled builder build
-> builder SBOM, signature, provenance, scan, policy, and approval evidence
-> immutable builder promotion
-> consumer-side evidence re-verification and local materialization
-> isolated consumer build without public-network fallback
```

## Consequences

- The Go builder has its own source, artifact, evidence, signer, promoter, and
  revocation boundaries.
- A Go consumer must reference a full approved builder digest and must not
  rebuild or substitute the builder in its own build, release, or hotfix lane.
- Builder caches, mutable tags, public registry fallbacks, and tenant evidence
  are prohibited as authority inputs.
- Go builder release, support, and hotfix branch families are introduced only
  with their complete governed lifecycle and corresponding immutable artifact
  delivery contract.
