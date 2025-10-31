package handlers

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/util"
	"github.com/sethvargo/go-githubactions"
)

// HandleReview processes pull_request_review events
// Automatically adds review/lgtm label when a PR is approved via GitHub UI by a collaborator
func HandleReview(ctx context.Context, client *github.Client, action *githubactions.Action, event *github.PullRequestReviewEvent) error {
	// Only handle approvals
	if event.GetReview().GetState() != util.ReviewStateApprovedWebhook {
		return nil
	}

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	prNum := event.GetPullRequest().GetNumber()
	reviewer := event.GetReview().GetUser().GetLogin()
	prAuthor := event.GetPullRequest().GetUser().GetLogin()

	// Prevent PR authors from approving their own changes
	if reviewer == prAuthor {
		action.Infof("PR #%d: ignoring self-approval from author %s", prNum, reviewer)
		return nil
	}

	// Check if reviewer is a collaborator
	isCollab, err := util.IsCollaborator(ctx, client, owner, repo, reviewer)
	if err != nil {
		return fmt.Errorf("failed to check collaborator status: %w", err)
	}

	if !isCollab {
		action.Infof("PR #%d approved by %s (non-collaborator) via GitHub UI, review recorded but label not added", prNum, reviewer)
		return nil
	}

	action.Infof("PR #%d approved by %s (collaborator) via GitHub UI, adding review/lgtm label", prNum, reviewer)

	// Add review/lgtm label and remove review/needs-review
	return util.ReplaceLabel(ctx, client, owner, repo, prNum, util.LabelReviewNeedsReview, util.LabelReviewLGTM)
}
