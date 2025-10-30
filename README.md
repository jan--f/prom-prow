# Prometheus Prow Bot

A lightweight chat-ops bot for Prometheus repositories, providing Prow-style commands without the complexity of running a full Prow instance.

## Features

- **Chat-ops commands** via PR/issue comments
- **Bidirectional /lgtm sync** - `/lgtm` command and GitHub UI approvals both add `review/lgtm` label
- **Automatic LGTM removal** when PR is updated with new commits
- **No OWNERS files needed** - uses GitHub's native collaborator permissions
- **Simple setup** - just a GitHub Action workflow

## Supported Commands

### `/lgtm`
Approves the PR and optionally adds the `review/lgtm` label.
- **Anyone**: Submits an approving GitHub review
- **Collaborators**: Also adds `review/lgtm` label and removes `review/needs-review` label
- **Note**: PR authors cannot approve their own changes

Example:
```
/lgtm
```

### `/lgtm cancel`
Cancels approval and optionally removes the `review/lgtm` label.
- **Anyone**: Dismisses the user's approval
- **Collaborators**: Also removes `review/lgtm` label (only if no other collaborators have approved) and adds `review/needs-review` label

Example:
```
/lgtm cancel
```

### `/cc @user1 @user2`
Requests reviews from specified users.
- Sends review requests via GitHub API
- **Requires**: Repository collaborator (write access)

Example:
```
/cc @roidelapluie @beorn7
```

### `/label <label>`
Adds one or more labels to the issue or PR.
- **Requires**: Repository collaborator (write access)

Examples:
```
/label component/promql
/label kind/bug priority/P1
```

### `/hold`
Adds the `blocked/hold` label to prevent merging.
- **Requires**: Repository collaborator (write access)

Example:
```
/hold
```

### `/unhold`
Removes the `blocked/hold` label.
- **Requires**: Repository collaborator (write access)

Example:
```
/unhold
```

## Automatic Behaviors

### Welcome Comment
When a new PR is opened:
- Posts a welcome comment with command instructions
- Only comments once (won't duplicate on workflow re-runs)

### GitHub UI Approval → Label
When a **collaborator** approves a PR via GitHub's UI (not using `/lgtm`):
- Automatically adds `review/lgtm` label
- Removes `review/needs-review` label
- Non-collaborator approvals are recorded but don't add labels
- PR authors cannot approve their own changes (self-approvals are ignored)

### New Commits → Remove LGTM
When new commits are pushed to a PR:
- Automatically removes `review/lgtm` label
- Adds `review/needs-review` label
- Forces re-review

## Installation

### For use in a Prometheus repository

1. Copy the bot to your repository as a local action or reference it as a reusable action
2. Create the workflow file `.github/workflows/prometheus-bot.yml`:

```yaml
name: Prometheus Bot

on:
  issue_comment:
    types: [created]
  pull_request_review:
    types: [submitted]
  pull_request:
    types: [opened, synchronize]

jobs:
  bot:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      issues: write
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - uses: prometheus/prom-prow@main
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

3. Ensure labels exist in your repository:
   - `review/needs-review`
   - `review/lgtm`
   - `blocked/hold`

## Development

### Building

```bash
make build
```

### Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

This generates `coverage.out` and `coverage.html` files.

### Docker

Build Docker image:
```bash
make docker-build
```

Run Docker image:
```bash
make docker-run
```

## Architecture

- **Language**: Go 1.21+
- **GitHub API**: `google/go-github/v57`
- **GitHub Actions**: `sethvargo/go-githubactions`
- **Deployment**: Docker container via GitHub Actions

## Comparison to Prow

| Feature | Prom-Prow Bot | Full Prow | Prow GitHub Actions |
|---------|---------------|-----------|---------------------|
| Chat-ops commands | ✅ | ✅ | ✅ |
| Welcome comment with instructions | ✅ | ✅ | ❌ |
| OWNERS files | ❌ (uses GitHub perms) | ✅ | ⚠️ (root only) |
| Auto review assignment | ❌ (use CODEOWNERS) | ✅ | ❌ |
| Bidirectional /lgtm | ✅ | ❌ | ❌ |
| Auto-remove LGTM on update | ✅ | ✅ | ✅ |
| Infrastructure | GitHub Actions | Kubernetes cluster | GitHub Actions |
| Setup complexity | Low | High | Medium |
| Maintenance burden | Low | High | Low |

## License

Apache 2.0
