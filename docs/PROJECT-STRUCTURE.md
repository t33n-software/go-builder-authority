# Project Structure

The repository separates builder-authority ownership from consumer and tenant
boundaries.

```text
.github/
  actions/                  reusable evidence and verification actions
  workflows/                source-quality and future builder delivery workflows
builder/
  go/                       Go builder definition and controlled input manifests
cmd/                        repository-local verification tools
docs/
  architecture/             authority and lifecycle decisions
  conventions/              Go builder source and hosting-platform rule-set
                            conventions
  development/              local and CI verification instructions
  operations/               evidence audit and revocation operations
  specification/            builder artifact and consumer contracts
internal/
  authority/                authority-owned Go domain and application logic
  packaging/                whitebox workflow and packaging contracts
policy/                     versioned policy references without credentials
```

The following boundaries remain external to this repository:

```text
tenant configuration
tenant secrets and identities
tenant application artifacts
tenant runtime and deployment evidence
credential-broker platform artifacts
```
