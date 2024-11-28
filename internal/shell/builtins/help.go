// internal/shell/builtins/help.go
package builtins

import (
	"fmt"

	"github.com/krzko/gosh/internal/shell/command"
)

type HelpCommand struct {
	commands map[string]command.BuiltinCommand
}

func NewHelpCommand(cmds map[string]command.BuiltinCommand) *HelpCommand {
	return &HelpCommand{
		commands: cmds,
	}
}

func (h *HelpCommand) Execute(args []string) error {
	if len(args) == 0 {
		fmt.Println("Available commands:")
		for name := range h.commands {
			fmt.Printf("  %s\n", name)
		}
		return nil
	}

	cmd, ok := h.commands[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %s", args[0])
	}
	fmt.Println(cmd.Help())
	return nil
}

func (h *HelpCommand) Help() string {
	return "help: Display help for commands\nUsage: help [command]"
}
