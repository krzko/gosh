// internal/shell/builtins/alias.go
package builtins

type AliasCommand struct {
	Cmd  BuiltinCommand
	Args []string
}

func (a *AliasCommand) Execute(args []string) error {
	// Combine the alias arguments with any additional arguments
	allArgs := append(a.Args, args...)
	return a.Cmd.Execute(allArgs)
}

func (a *AliasCommand) Help() string {
	return a.Cmd.Help()
}
