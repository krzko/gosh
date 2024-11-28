// internal/shell/builtins/exit.go
package builtins

import (
	"os"
)

type ExitCommand struct{}

func (e *ExitCommand) Execute(args []string) error {
	os.Exit(0)
	return nil
}

func (e *ExitCommand) Help() string {
	return "exit: Exit the shell\nUsage: exit"
}
