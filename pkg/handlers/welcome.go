package handlers

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

const welcomeComment = `Thank you for your pull request!

## Available Commands

You can use the following commands in comments to interact with this PR:

### Review Commands
- **/lgtm** - Approve this PR
  - Anyone: Submits an approving GitHub review
  - Collaborators: Also adds the review/lgtm label
  - Note: PR authors cannot approve their own changes
- **/lgtm cancel** - Cancel your approval
  - Anyone: Dismisses your review
  - Collaborators: Also removes the review/lgtm label (if no other collaborators have approved)

### Collaboration Commands *(require write access)*
- **/cc @user1 @user2** - Request reviews from specific users
- **/label <name>** - Add label(s) to this PR (e.g., /label kind/bug component/promql)
- **/hold** - Mark this PR as on-hold (adds blocked/hold label to prevent merging)
- **/unhold** - Remove the hold (removes blocked/hold label)

## Automatic Behaviors

- When a collaborator approves via GitHub's UI, the bot automatically adds the review/lgtm label
- When new commits are pushed, the bot automatically removes review/lgtm to require re-review

## Permission Model

- **Anyone** can use /lgtm and /lgtm cancel to submit reviews
- **Collaborators** (write access) can additionally manage labels via /lgtm, use /cc, /label, /hold, and /unhold

This allows community members to provide reviews while maintaining label control for maintainers.

---
*This is an automated message from the Prometheus Prow Bot. For issues or questions, please file an issue in the prometheus/prom-prow repository.*`

// HandleWelcome posts a welcome comment on newly opened PRs
func HandleWelcome(ctx context.Context, client *github.Client, action *githubactions.Action, event *github.PullRequestEvent) error {
	// Only handle opened PRs
	if event.GetAction() != "opened" {
		return nil
	}

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	prNum := event.GetPullRequest().GetNumber()

	action.Infof("Posting welcome comment on new PR #%d", prNum)

	// Check if we already commented (in case of re-runs)
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNum, nil)
	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	// Get bot user (the GitHub Actions bot)
	botLogin := "github-actions[bot]"

	// Check if bot already commented
	for _, comment := range comments {
		if comment.GetUser().GetLogin() == botLogin {
			action.Infof("Bot already commented on PR #%d, skipping welcome message", prNum)
			return nil
		}
	}

	// Post welcome comment
	comment := &github.IssueComment{
		Body: github.String(welcomeComment),
	}

	_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNum, comment)
	if err != nil {
		return fmt.Errorf("failed to create welcome comment: %w", err)
	}

	return nil
}
