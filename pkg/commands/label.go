package commands

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/util"
)

// HandleLabel handles /label command to add labels
func HandleLabel(ctx context.Context, client *github.Client, owner, repo string, number int, labels []string) error {
	if len(labels) == 0 {
		return fmt.Errorf("no labels specified for /label command")
	}

	for _, label := range labels {
		if err := util.AddLabel(ctx, client, owner, repo, number, label); err != nil {
			return err
		}
	}

	return nil
}
