package handlers

import (
	"context"
	"testing"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

func TestHandlePRUpdate_NonSynchronize(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)
	action := githubactions.New()

	tests := []struct {
		name   string
		action string
	}{
		{
			name:   "opened",
			action: "opened",
		},
		{
			name:   "closed",
			action: "closed",
		},
		{
			name:   "reopened",
			action: "reopened",
		},
		{
			name:   "edited",
			action: "edited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &github.PullRequestEvent{
				Action: github.String(tt.action),
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

			// Should return nil for non-synchronize actions
			err := HandlePRUpdate(ctx, client, action, event)
			if err != nil {
				t.Errorf("HandlePRUpdate() should ignore %s action, got error: %v", tt.action, err)
			}
		})
	}
}
