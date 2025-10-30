package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/handlers"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	// Get GitHub Actions context
	action := githubactions.New()

	// Get inputs
	token := action.GetInput("github-token")
	if token == "" {
		action.Fatalf("github-token is required")
	}

	// Create GitHub client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get event information
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	if eventPath == "" {
		action.Fatalf("GITHUB_EVENT_PATH not set")
	}

	// Read event payload
	eventData, err := os.ReadFile(eventPath)
	if err != nil {
		action.Fatalf("Failed to read event file: %v", err)
	}

	// Route to appropriate handler based on event type
	if err := handleEvent(ctx, client, action, eventName, eventData); err != nil {
		action.Fatalf("Failed to handle event: %v", err)
	}

	action.Infof("Successfully processed %s event", eventName)
}

func handleEvent(ctx context.Context, client *github.Client, action *githubactions.Action, eventName string, eventData []byte) error {
	switch eventName {
	case "issue_comment":
		var event github.IssueCommentEvent
		if err := json.Unmarshal(eventData, &event); err != nil {
			return fmt.Errorf("failed to unmarshal issue_comment event: %w", err)
		}
		return handlers.HandleComment(ctx, client, action, &event)

	case "pull_request_review":
		var event github.PullRequestReviewEvent
		if err := json.Unmarshal(eventData, &event); err != nil {
			return fmt.Errorf("failed to unmarshal pull_request_review event: %w", err)
		}
		return handlers.HandleReview(ctx, client, action, &event)

	case "pull_request":
		var event github.PullRequestEvent
		if err := json.Unmarshal(eventData, &event); err != nil {
			return fmt.Errorf("failed to unmarshal pull_request event: %w", err)
		}

		// Handle welcome comment for opened PRs
		if err := handlers.HandleWelcome(ctx, client, action, &event); err != nil {
			action.Errorf("Failed to handle welcome: %v", err)
		}

		// Handle PR updates (synchronize)
		return handlers.HandlePRUpdate(ctx, client, action, &event)

	default:
		action.Warningf("Unsupported event type: %s", eventName)
		return nil
	}
}
