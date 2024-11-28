// internal/shell/command/command.go
package command

type BuiltinCommand interface {
	Execute(args []string) error
	Help() string
}
