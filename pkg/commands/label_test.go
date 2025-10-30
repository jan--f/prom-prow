package commands

import (
	"context"
	"testing"

	"github.com/google/go-github/v57/github"
)

func TestHandleLabel_EmptyList(t *testing.T) {
	ctx := context.Background()
	client := github.NewClient(nil)

	// Test with empty labels list - should return error before API call
	err := HandleLabel(ctx, client, "owner", "repo", 1, []string{})
	if err == nil {
		t.Error("Expected error for empty labels, got nil")
	}
	if err.Error() != "no labels specified for /label command" {
		t.Errorf("Expected 'no labels specified' error, got: %v", err)
	}
}
