package handlers

import (
	"context"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/util"
	"github.com/sethvargo/go-githubactions"
)

// HandlePRUpdate processes pull_request synchronize events
// Automatically removes review/lgtm label when new commits are pushed
func HandlePRUpdate(ctx context.Context, client *github.Client, action *githubactions.Action, event *github.PullRequestEvent) error {
	// Only handle synchronize (new commits)
	if event.GetAction() != "synchronize" {
		return nil
	}

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	prNum := event.GetPullRequest().GetNumber()

	action.Infof("PR #%d updated with new commits, removing review/lgtm label", prNum)

	// Remove review/lgtm and add review/needs-review
	return util.ReplaceLabel(ctx, client, owner, repo, prNum, util.LabelReviewLGTM, util.LabelReviewNeedsReview)
}
