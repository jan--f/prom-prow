package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
)

// HandleCC handles /cc command to request reviews
func HandleCC(ctx context.Context, client *github.Client, owner, repo string, prNum int, users []string) error {
	if len(users) == 0 {
		return fmt.Errorf("no users specified for /cc command")
	}

	// Remove @ prefix if present
	reviewers := make([]string, 0, len(users))
	for _, user := range users {
		user = strings.TrimPrefix(user, "@")
		if user != "" {
			reviewers = append(reviewers, user)
		}
	}

	if len(reviewers) == 0 {
		return fmt.Errorf("no valid users specified for /cc command")
	}

	// Request reviews
	reviewRequest := github.ReviewersRequest{
		Reviewers: reviewers,
	}

	_, _, err := client.PullRequests.RequestReviewers(ctx, owner, repo, prNum, reviewRequest)
	if err != nil {
		return fmt.Errorf("failed to request reviewers: %w", err)
	}

	return nil
}
