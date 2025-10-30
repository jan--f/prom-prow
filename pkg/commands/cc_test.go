package commands

import (
	"context"
	"testing"

	"github.com/google/go-github/v57/github"
)

func TestHandleCC_EmptyList(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)

	// Test with empty users list - should return error before API call
	err := HandleCC(ctx, client, "owner", "repo", 1, []string{})
	if err == nil {
		t.Error("Expected error for empty users, got nil")
	}
	if err.Error() != "no users specified for /cc command" {
		t.Errorf("Expected 'no users specified' error, got: %v", err)
	}
}

func TestHandleCC_EmptyAfterTrimming(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)

	// Test with @ symbols only - should return error before API call
	err := HandleCC(ctx, client, "owner", "repo", 1, []string{"@", "@"})
	if err == nil {
		t.Error("Expected error for invalid users, got nil")
	}
	if err.Error() != "no valid users specified for /cc command" {
		t.Errorf("Expected 'no valid users' error, got: %v", err)
	}
}
