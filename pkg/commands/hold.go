package commands

import (
	"context"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/util"
)

// HandleHold handles /hold command to add blocked/hold label
func HandleHold(ctx context.Context, client *github.Client, owner, repo string, number int) error {
	return util.AddLabel(ctx, client, owner, repo, number, "blocked/hold")
}

// HandleUnhold handles /unhold command to remove blocked/hold label
func HandleUnhold(ctx context.Context, client *github.Client, owner, repo string, number int) error {
	return util.RemoveLabel(ctx, client, owner, repo, number, "blocked/hold")
}
