package util

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantCmd  string
		wantArgs []string
		wantNil  bool
	}{
		{
			name:     "simple lgtm",
			body:     "/lgtm",
			wantCmd:  "lgtm",
			wantArgs: []string{},
		},
		{
			name:     "lgtm cancel",
			body:     "/lgtm cancel",
			wantCmd:  "lgtm",
			wantArgs: []string{"cancel"},
		},
		{
			name:     "cc with users",
			body:     "/cc @user1 @user2",
			wantCmd:  "cc",
			wantArgs: []string{"@user1", "@user2"},
		},
		{
			name:     "label with multiple labels",
			body:     "/label kind/bug component/promql",
			wantCmd:  "label",
			wantArgs: []string{"kind/bug", "component/promql"},
		},
		{
			name:     "hold",
			body:     "/hold",
			wantCmd:  "hold",
			wantArgs: []string{},
		},
		{
			name:     "unhold",
			body:     "/unhold",
			wantCmd:  "unhold",
			wantArgs: []string{},
		},
		{
			name:     "command in multiline comment",
			body:     "This looks good\n/lgtm\nThanks!",
			wantCmd:  "lgtm",
			wantArgs: []string{},
		},
		{
			name:     "command with leading whitespace",
			body:     "  /lgtm  ",
			wantCmd:  "lgtm",
			wantArgs: []string{},
		},
		{
			name:    "no command",
			body:    "This is just a comment",
			wantNil: true,
		},
		{
			name:    "empty body",
			body:    "",
			wantNil: true,
		},
		{
			name:    "invalid command format",
			body:    "lgtm without slash",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ParseCommand(tt.body)
			if tt.wantNil {
				if cmd != nil {
					t.Errorf("ParseCommand() = %v, want nil", cmd)
				}
				return
			}

			if cmd == nil {
				t.Fatalf("ParseCommand() = nil, want command")
			}

			if cmd.Name != tt.wantCmd {
				t.Errorf("ParseCommand().Name = %v, want %v", cmd.Name, tt.wantCmd)
			}

			// Handle nil vs empty slice comparison
			gotArgs := cmd.Args
			if gotArgs == nil {
				gotArgs = []string{}
			}
			wantArgs := tt.wantArgs
			if wantArgs == nil {
				wantArgs = []string{}
			}

			if len(gotArgs) != len(wantArgs) {
				t.Errorf("ParseCommand().Args = %v (len=%d), want %v (len=%d)", gotArgs, len(gotArgs), wantArgs, len(wantArgs))
				return
			}

			for i, arg := range gotArgs {
				if arg != wantArgs[i] {
					t.Errorf("ParseCommand().Args[%d] = %v, want %v", i, arg, wantArgs[i])
				}
			}
		})
	}
}
