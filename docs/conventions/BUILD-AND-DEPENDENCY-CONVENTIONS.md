# Build and Dependency Conventions

The GBA-1 source foundation uses Go 1.26.5 with:

```text
GOTOOLCHAIN=local
GOFLAGS=-mod=readonly
GOVCS=*:off
```

GBA-1 establishes source-quality controls only. It does not treat the current
public Go module proxy, a Docker cache, or a runner cache as an approved
builder or dependency authority.

Before a production builder artifact can be released, the authority requires:

```text
an approved internal Go proxy
controlled dependency resolution and admission
immutable approved base-image inputs
an independently verified builder artifact
```
