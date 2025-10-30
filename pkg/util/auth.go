package util

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
)

// IsCollaborator checks if a user has write access to the repository
func IsCollaborator(ctx context.Context, client *github.Client, owner, repo, user string) (bool, error) {
	perm, _, err := client.Repositories.GetPermissionLevel(ctx, owner, repo, user)
	if err != nil {
		return false, fmt.Errorf("failed to get permission level: %w", err)
	}

	// Users with write, maintain, or admin access are considered collaborators
	permission := perm.GetPermission()
	return permission == "write" || permission == "maintain" || permission == "admin", nil
}
