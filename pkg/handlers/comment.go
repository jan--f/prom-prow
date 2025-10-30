package handlers

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/prom-prow/pkg/commands"
	"github.com/prometheus/prom-prow/pkg/util"
	"github.com/sethvargo/go-githubactions"
)

// HandleComment processes issue_comment events
func HandleComment(ctx context.Context, client *github.Client, action *githubactions.Action, event *github.IssueCommentEvent) error {
	// Only handle created comments
	if event.GetAction() != "created" {
		return nil
	}

	// Only handle PR comments
	if !event.GetIssue().IsPullRequest() {
		return nil
	}

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	prNum := event.GetIssue().GetNumber()
	commenter := event.GetComment().GetUser().GetLogin()
	body := event.GetComment().GetBody()

	// Parse command
	cmd := util.ParseCommand(body)
	if cmd == nil {
		// No command found, ignore
		return nil
	}

	action.Infof("Processing command /%s from %s on PR #%d", cmd.Name, commenter, prNum)

	// Check if user is a collaborator (write access)
	isCollab, err := util.IsCollaborator(ctx, client, owner, repo, commenter)
	if err != nil {
		return fmt.Errorf("failed to check collaborator status: %w", err)
	}

	// Commands that require write access
	requiresWriteAccess := map[string]bool{
		"cc":     true,
		"hold":   true,
		"unhold": true,
		"label":  true,
	}

	if requiresWriteAccess[cmd.Name] && !isCollab {
		action.Warningf("User %s is not a collaborator, ignoring /%s command (requires write access)", commenter, cmd.Name)
		return nil
	}

	// Route to command handler
	switch cmd.Name {
	case "lgtm":
		cancel := len(cmd.Args) > 0 && cmd.Args[0] == "cancel"
		return commands.HandleLGTM(ctx, client, owner, repo, prNum, commenter, cancel, isCollab)

	case "cc":
		return commands.HandleCC(ctx, client, owner, repo, prNum, cmd.Args)

	case "label":
		return commands.HandleLabel(ctx, client, owner, repo, prNum, cmd.Args)

	case "hold":
		return commands.HandleHold(ctx, client, owner, repo, prNum)

	case "unhold":
		return commands.HandleUnhold(ctx, client, owner, repo, prNum)

	default:
		action.Warningf("Unknown command: /%s", cmd.Name)
		return nil
	}
}
