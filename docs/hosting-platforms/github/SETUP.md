# GitHub Setup

All GitHub repository configuration is performed through the GitHub graphical
interface.

## Repository settings

Open:

```text
t33n-software/go-builder-authority
-> Settings
-> General
```

Set:

```text
Allow merge commits: enabled
Allow rebase merging: enabled
Allow squash merging: enabled
Automatically delete head branches: enabled
Allow auto-merge: disabled
Always suggest updating pull request branches: disabled
Enable release immutability: enabled
```

## Source publisher identity

Before publishing the first GBA-1 working branch, create a dedicated GitHub
App:

```text
Name:
go-builder-authority-source-publisher

Repository access:
Only select repositories

Repository:
t33n-software/go-builder-authority

Webhook:
disabled

Repository permissions:
Contents: Read and write
Pull requests: Read and write
Metadata: Read-only automatically
```

Do not grant:

```text
Actions
Workflows
Administration
Secrets
Environments
Code scanning
Ruleset bypass
Organization permissions
```

The publisher private key belongs only in a dedicated platform secret boundary.
It must not be stored on a developer machine, in GitHub variables, in Git, or
in this repository.

## Ruleset import

Do not import any Ruleset before the first GBA-1 pull request has emitted the
actual Quality, CodeQL, and dependency-admission results.

After that pull request succeeds, import the three JSON files from
`rulesets/` in numeric order. `release/*` and `support/*` Rulesets remain
absent until the authority has its own governed builder artifact release and
maintenance lifecycle.
