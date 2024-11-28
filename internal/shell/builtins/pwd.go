// internal/shell/builtins/pwd.go
package builtins

import (
	"fmt"
	"os"
)

type PwdCommand struct{}

func (p *PwdCommand) Execute(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func (p *PwdCommand) Help() string {
	return "pwd: Print working directory\nUsage: pwd"
}
