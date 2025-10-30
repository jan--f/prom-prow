package commands

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/util"
)

// HandleLGTM handles /lgtm and /lgtm cancel commands
// isCollaborator determines if labels should be added/removed
func HandleLGTM(ctx context.Context, client *github.Client, owner, repo string, prNum int, user, prAuthor string, cancel bool, isCollaborator bool) error {
	// Prevent PR authors from approving their own changes
	if user == prAuthor {
		return fmt.Errorf("PR authors cannot approve their own changes")
	}

	if cancel {
		return cancelLGTM(ctx, client, owner, repo, prNum, user, isCollaborator)
	}

	return approveLGTM(ctx, client, owner, repo, prNum, user, isCollaborator)
}

func approveLGTM(ctx context.Context, client *github.Client, owner, repo string, prNum int, user string, isCollaborator bool) error {
	// Submit approving review (anyone can do this)
	review := &github.PullRequestReviewRequest{
		Event: github.String("APPROVE"),
		Body:  github.String("LGTM"),
	}

	_, _, err := client.PullRequests.CreateReview(ctx, owner, repo, prNum, review)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	// Only collaborators can add/remove labels
	if isCollaborator {
		if err := util.ReplaceLabel(ctx, client, owner, repo, prNum, "review/needs-review", "review/lgtm"); err != nil {
			return err
		}
	}

	return nil
}

func cancelLGTM(ctx context.Context, client *github.Client, owner, repo string, prNum int, user string, isCollaborator bool) error {
	// Get all reviews to find the user's review
	reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, prNum, nil)
	if err != nil {
		return fmt.Errorf("failed to list reviews: %w", err)
	}

	// Find and dismiss the user's approval (anyone can dismiss their own review)
	for _, review := range reviews {
		if review.GetUser().GetLogin() == user && review.GetState() == "APPROVED" {
			dismissal := &github.PullRequestReviewDismissalRequest{
				Message: github.String("LGTM cancelled"),
			}
			_, _, err := client.PullRequests.DismissReview(ctx, owner, repo, prNum, review.GetID(), dismissal)
			if err != nil {
				return fmt.Errorf("failed to dismiss review: %w", err)
			}
			break
		}
	}

	// Only collaborators can remove the label
	// But only remove it if there are no remaining approvals from collaborators
	if isCollaborator {
		// Check if any other collaborators still have approvals
		hasCollaboratorApproval, err := hasRemainingCollaboratorApprovals(ctx, client, owner, repo, prNum, user)
		if err != nil {
			return err
		}

		// Only remove label if no collaborator approvals remain
		if !hasCollaboratorApproval {
			if err := util.ReplaceLabel(ctx, client, owner, repo, prNum, "review/lgtm", "review/needs-review"); err != nil {
				return err
			}
		}
	}

	return nil
}

// hasRemainingCollaboratorApprovals checks if there are any approved reviews from collaborators
// excluding the specified user
func hasRemainingCollaboratorApprovals(ctx context.Context, client *github.Client, owner, repo string, prNum int, excludeUser string) (bool, error) {
	reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, prNum, nil)
	if err != nil {
		return false, fmt.Errorf("failed to list reviews: %w", err)
	}

	// Track the latest review state for each user (excluding the user who cancelled)
	latestReviewState := make(map[string]string)
	for _, review := range reviews {
		reviewer := review.GetUser().GetLogin()
		if reviewer == excludeUser {
			continue
		}
		// GitHub returns reviews in chronological order, so later reviews override earlier ones
		state := review.GetState()
		if state != "" {
			latestReviewState[reviewer] = state
		}
	}

	// Check if any of the remaining reviewers with APPROVED state are collaborators
	for reviewer, state := range latestReviewState {
		if state == "APPROVED" {
			isCollab, err := util.IsCollaborator(ctx, client, owner, repo, reviewer)
			if err != nil {
				return false, fmt.Errorf("failed to check collaborator status for %s: %w", reviewer, err)
			}
			if isCollab {
				return true, nil
			}
		}
	}

	return false, nil
}
