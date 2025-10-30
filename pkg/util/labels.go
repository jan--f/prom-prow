package util

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
)

// AddLabel adds a label to an issue or PR, ignoring if it already exists
func AddLabel(ctx context.Context, client *github.Client, owner, repo string, number int, label string) error {
	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, number, []string{label})
	if err != nil {
		return fmt.Errorf("failed to add label %s: %w", label, err)
	}
	return nil
}

// RemoveLabel removes a label from an issue or PR, ignoring if it doesn't exist
func RemoveLabel(ctx context.Context, client *github.Client, owner, repo string, number int, label string) error {
	_, err := client.Issues.RemoveLabelForIssue(ctx, owner, repo, number, label)
	if err != nil {
		// GitHub returns 404 if label doesn't exist, which is fine
		if ghErr, ok := err.(*github.ErrorResponse); ok && ghErr.Response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("failed to remove label %s: %w", label, err)
	}
	return nil
}

// ReplaceLabel removes oldLabel and adds newLabel
func ReplaceLabel(ctx context.Context, client *github.Client, owner, repo string, number int, oldLabel, newLabel string) error {
	if err := RemoveLabel(ctx, client, owner, repo, number, oldLabel); err != nil {
		return err
	}
	return AddLabel(ctx, client, owner, repo, number, newLabel)
}
