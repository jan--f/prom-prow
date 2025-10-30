package handlers

import (
	"context"
	"testing"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

func TestHandleComment_NonPRComment(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	action := githubactions.New()

	// Create event for issue (not PR) comment
	event := &github.IssueCommentEvent{
		Action: github.String("created"),
		Issue: &github.Issue{
			Number:        github.Int(1),
			PullRequestLinks: nil, // Not a PR
		},
		Comment: &github.IssueComment{
			Body: github.String("/lgtm"),
			User: &github.User{
				Login: github.String("user"),
			},
		},
		Repo: &github.Repository{
			Name: github.String("repo"),
			Owner: &github.User{
				Login: github.String("owner"),
			},
		},
	}

	// Should return nil for non-PR comments
	err := HandleComment(ctx, client, action, event)
	if err != nil {
		t.Errorf("HandleComment() should ignore non-PR comments, got error: %v", err)
	}
}

func TestHandleComment_NonCreatedAction(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	action := githubactions.New()

	// Create event for edited comment
	event := &github.IssueCommentEvent{
		Action: github.String("edited"),
		Issue: &github.Issue{
			Number: github.Int(1),
			PullRequestLinks: &github.PullRequestLinks{
				URL: github.String("url"),
			},
		},
		Comment: &github.IssueComment{
			Body: github.String("/lgtm"),
			User: &github.User{
				Login: github.String("user"),
			},
		},
		Repo: &github.Repository{
			Name: github.String("repo"),
			Owner: &github.User{
				Login: github.String("owner"),
			},
		},
	}

	// Should return nil for non-created actions
	err := HandleComment(ctx, client, action, event)
	if err != nil {
		t.Errorf("HandleComment() should ignore non-created actions, got error: %v", err)
	}
}

func TestHandleComment_NoCommand(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	action := githubactions.New()

	// Create event with no command
	event := &github.IssueCommentEvent{
		Action: github.String("created"),
		Issue: &github.Issue{
			Number: github.Int(1),
			PullRequestLinks: &github.PullRequestLinks{
				URL: github.String("url"),
			},
		},
		Comment: &github.IssueComment{
			Body: github.String("Just a regular comment"),
			User: &github.User{
				Login: github.String("user"),
			},
		},
		Repo: &github.Repository{
			Name: github.String("repo"),
			Owner: &github.User{
				Login: github.String("owner"),
			},
		},
	}

	// Should return nil when no command found
	err := HandleComment(ctx, client, action, event)
	if err != nil {
		t.Errorf("HandleComment() should ignore comments without commands, got error: %v", err)
	}
}
