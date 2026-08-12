# GitHub Rulesets

These files map the canonical Git governance to the current Go Builder
Authority source boundary. They do not define a second release, evidence, or
builder policy.

## Import timing

Import through the GitHub graphical interface only after the first GBA-1 pull
request has produced the actual required checks for this repository:

```text
Settings
-> Rules
-> Rulesets
-> New ruleset
-> Import a ruleset
```

Import in this order:

```text
01-ticket-working-branches.json
02-develop.json
03-main.json
```

Do not import a `release/*` or `support/*` Ruleset yet. Those families require
their own governed builder-release workflow, immutable builder artifact
delivery, and release or maintenance contract.

## Required GitHub repository settings

```text
Allow merge commits: enabled
Allow rebase merging: enabled
Allow squash merging: enabled
Automatically delete head branches: enabled
Allow auto-merge: disabled
Always suggest updating pull request branches: disabled
Enable release immutability: enabled
```

## Required checks

The shared-line Rulesets require only checks emitted by the current source
workflows:

```text
Quality gates (linux-amd64)
Dependency admission review
CodeQL code scanning with all alerts blocking
```

The Linux quality gate enforces formatting, module integrity, tests, exact
100% statement coverage, race detection, static analysis, Linux/AMD64 build,
and module provenance.

## Security boundary

Ruleset files must not contain builder registry credentials, signing keys,
tenant configuration, tenant identities, bypass actors, or mutable builder
references.
