// internal/shell/shell.go
package shell

import (
	"fmt"
	"os"

	"github.com/krzko/gosh/internal/shell/completion"
	"github.com/krzko/gosh/internal/shell/executor"
	"github.com/krzko/gosh/internal/shell/history"
	"github.com/krzko/gosh/internal/shell/parser"
	"github.com/krzko/gosh/internal/shell/prompt"
)

type Shell struct {
	prompt    *prompt.Manager
	parser    *parser.Parser
	executor  *executor.Executor
	completer *completion.Completer
	history   *history.Manager
}

// internal/shell/shell.go
func New() (*Shell, error) {
	// Initialize history first
	hist, err := history.NewManager(".gosh_history")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize history: %w", err)
	}

	// Initialize executor to get builtins
	executor := executor.New()

	// Initialize completer
	completer := completion.NewCompleter()

	// Initialize prompt with history and builtins
	promptManager, err := prompt.NewManager(completer, executor.GetBuiltins())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prompt: %w", err)
	}

	// Create shell instance
	sh := &Shell{
		prompt:    promptManager,
		parser:    parser.New(),
		executor:  executor,
		completer: completer,
		history:   hist,
	}

	return sh, nil
}

func (s *Shell) Run() error {
	defer s.cleanup()

	for {
		input, err := s.prompt.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return nil // Clean exit on Ctrl+D
			}
			return err
		}

		if input == "" {
			continue
		}

		// Add to history
		if err := s.history.Add(input); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add to history: %v\n", err)
		}

		// Parse the command
		cmd, err := s.parser.Parse(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			continue
		}

		// Execute the command
		if err := s.executor.Execute(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		}
	}
}

func (s *Shell) cleanup() {
	if s.prompt != nil {
		s.prompt.Close()
	}
	if s.history != nil {
		s.history.Close()
	}
}
