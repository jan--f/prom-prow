package handlers

import (
	"context"
	"testing"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

func TestHandleReview_NonApproval(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	action := githubactions.New()

	// Create event for non-approval review
	event := &github.PullRequestReviewEvent{
		Review: &github.PullRequestReview{
			State: github.String("commented"),
			User: &github.User{
				Login: github.String("user"),
			},
		},
		PullRequest: &github.PullRequest{
			Number: github.Int(1),
		},
		Repo: &github.Repository{
			Name: github.String("repo"),
			Owner: &github.User{
				Login: github.String("owner"),
			},
		},
	}

	// Should return nil for non-approval reviews
	err := HandleReview(ctx, client, action, event)
	if err != nil {
		t.Errorf("HandleReview() should ignore non-approval reviews, got error: %v", err)
	}
}

func TestHandleReview_ChangeRequested(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	action := githubactions.New()

	// Create event for changes requested
	event := &github.PullRequestReviewEvent{
		Review: &github.PullRequestReview{
			State: github.String("changes_requested"),
			User: &github.User{
				Login: github.String("user"),
			},
		},
		PullRequest: &github.PullRequest{
			Number: github.Int(1),
		},
		Repo: &github.Repository{
			Name: github.String("repo"),
			Owner: &github.User{
				Login: github.String("owner"),
			},
		},
	}

	// Should return nil for changes_requested
	err := HandleReview(ctx, client, action, event)
	if err != nil {
		t.Errorf("HandleReview() should ignore changes_requested, got error: %v", err)
	}
}
