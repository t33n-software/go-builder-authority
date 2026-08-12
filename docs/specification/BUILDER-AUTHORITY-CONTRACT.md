# Go Builder Authority Contract

The authority produces approved Go builder artifacts, not application or tenant
artifacts.

## Required builder subject

Before a consumer can start an isolated build, the authority must publish an
immutable builder subject containing:

```text
builder source and pinned base-image inputs
exact Go toolchain version
full builder image digest
SBOM
signature
provenance and attestation
scan, policy, approval, and revocation evidence
immutable artifact and evidence references
```

## Consumer contract

Consumers must:

1. obtain the builder only through an approved internal registry and full
   immutable digest;
2. re-verify all required builder evidence;
3. materialize the exact builder locally;
4. verify local digest equality;
5. run the consumer build without public-network fallback or runtime pulls.

Consumers must not rebuild, retag, substitute, or use cached copies of the
builder as a trust authority.

## Authority exclusions

The authority must not receive tenant secrets, tenant runtime identities,
tenant deployment authority, tenant application artifacts, or tenant evidence.
