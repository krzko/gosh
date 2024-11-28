// internal/shell/executor/executor.go
package executor

import (
	"os"
	"os/exec"

	"github.com/krzko/gosh/internal/shell/builtins"
	"github.com/krzko/gosh/internal/shell/command"
)

type Command struct {
	Name string
	Args []string
	Pipe *Command
}

type Executor struct {
	builtins map[string]command.BuiltinCommand
}

func New() *Executor {
	return &Executor{
		builtins: registerBuiltins(),
	}
}

func (e *Executor) Execute(cmd *Command) error {
	if builtin, ok := e.builtins[cmd.Name]; ok {
		return builtin.Execute(cmd.Args)
	}

	return e.executeExternal(cmd)
}

func (e *Executor) executeExternal(cmd *Command) error {
	if cmd.Pipe != nil {
		return e.executePipeline(cmd)
	}

	command := exec.Command(cmd.Name, cmd.Args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func (e *Executor) executePipeline(first *Command) error {
	var commands []*exec.Cmd
	current := first

	for current != nil {
		cmd := exec.Command(current.Name, current.Args...)
		commands = append(commands, cmd)
		current = current.Pipe
	}

	// Connect the pipeline
	for i := 0; i < len(commands)-1; i++ {
		pipe, err := commands[i].StdoutPipe()
		if err != nil {
			return err
		}
		commands[i+1].Stdin = pipe
	}

	// Start all commands
	for _, cmd := range commands {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	// Wait for all commands
	for _, cmd := range commands {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) GetBuiltins() map[string]command.BuiltinCommand {
	return e.builtins
}

func registerBuiltins() map[string]command.BuiltinCommand {
	builtinMap := map[string]command.BuiltinCommand{
		"ls":    builtins.NewLsCommand(),
		"cd":    &builtins.CdCommand{},
		"http":  &builtins.HttpCommand{},
		"https": &builtins.HttpsCommand{},
		"exit":  &builtins.ExitCommand{},
		"pwd":   &builtins.PwdCommand{},
		"ver":   &builtins.VerCommand{},
	}

	// Aliases
	builtinMap["ll"] = &builtins.AliasCommand{
		Cmd:  builtinMap["ls"],
		Args: []string{"-la"},
	}

	builtinMap["help"] = builtins.NewHelpCommand(builtinMap)

	return builtinMap
}
