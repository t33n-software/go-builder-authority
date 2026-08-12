# Go Builder Definition

This directory will contain the controlled Go builder definition and immutable
input manifests.

GBA-1 deliberately does not add a Dockerfile or a public Go base-image
reference. A builder definition is introduced only after the authority binds:

```text
approved base-image digest
exact Go toolchain and tool digests
approved internal Go proxy
builder artifact and evidence registry references
builder signer, promoter, policy, approval, and revocation boundaries
```

Until then, no consumer may treat this directory, a runner cache, or a public
image as an approved builder artifact.
