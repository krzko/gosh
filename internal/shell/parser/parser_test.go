// internal/shell/parser/parser_test.go
package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCmd  string
		wantArgs []string
		wantPipe bool
	}{
		{
			name:     "simple command",
			input:    "ls -l",
			wantCmd:  "ls",
			wantArgs: []string{"-l"},
			wantPipe: false,
		},
		{
			name:     "command with quotes",
			input:    `echo "hello world"`,
			wantCmd:  "echo",
			wantArgs: []string{"hello world"},
			wantPipe: false,
		},
		{
			name:     "command with pipe",
			input:    "ls -l | grep foo",
			wantCmd:  "ls",
			wantArgs: []string{"-l"},
			wantPipe: true,
		},
		{
			name:     "command with escaped characters",
			input:    `echo hello\ world`,
			wantCmd:  "echo",
			wantArgs: []string{"hello world"},
			wantPipe: false,
		},
	}

	parser := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if cmd.Name != tt.wantCmd {
				t.Errorf("Parse() command = %v, want %v", cmd.Name, tt.wantCmd)
			}

			if len(cmd.Args) != len(tt.wantArgs) {
				t.Errorf("Parse() args length = %v, want %v", len(cmd.Args), len(tt.wantArgs))
			}

			for i, arg := range cmd.Args {
				if arg != tt.wantArgs[i] {
					t.Errorf("Parse() arg[%d] = %v, want %v", i, arg, tt.wantArgs[i])
				}
			}

			hasPipe := cmd.Pipe != nil
			if hasPipe != tt.wantPipe {
				t.Errorf("Parse() pipe = %v, want %v", hasPipe, tt.wantPipe)
			}
		})
	}
}
