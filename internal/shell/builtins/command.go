// internal/shell/builtins/command.go
package builtins

type BuiltinCommand interface {
	Execute(args []string) error
	Help() string
}
