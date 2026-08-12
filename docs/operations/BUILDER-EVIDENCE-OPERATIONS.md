# Builder Evidence Operations

Builder evidence is append-only and subject-bound. It is not replaced by CI
logs, container tags, runner caches, or visible workflow success.

This authority stores no tenant application, tenant runtime, or tenant
deployment evidence.

## Required operational evidence

```text
Source
-> Dependency Resolution
-> Build
-> Builder Artifact
-> Promotion
-> Deployment or Consumer Admission
-> Operation
```

Every subject must remain bound to immutable references and digests. A
revoked, quarantined, failed, pending, superseded, or not-recorded builder
subject authorizes no consumer build, artifact push, promotion, or deployment.

## Incident response

On builder compromise, revocation, or digest mismatch:

1. block dependent consumer lanes;
2. create an append-only operation subject for the affected builder and
   blast radius;
3. publish a separately verified replacement builder;
4. rebuild affected consumer artifacts through their own governed lifecycle.

Neither a local cache nor a mutable replacement tag is a recovery mechanism.
