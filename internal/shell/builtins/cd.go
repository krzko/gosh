// internal/shell/builtins/cd.go
package builtins

import (
	"os"
	"path/filepath"
)

type CdCommand struct{}

func (c *CdCommand) Execute(args []string) error {
	var dir string
	if len(args) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir = homeDir
	} else {
		dir = args[0]

		// Handle ".." special case
		if dir == ".." {
			currentDir, err := os.Getwd()
			if err != nil {
				return err
			}
			dir = filepath.Dir(currentDir)
		}
	}

	// Clean the path to resolve any .. or . in the middle of paths
	dir = filepath.Clean(dir)

	// Expand ~ to home directory if it's the first character
	if len(dir) > 0 && dir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir = filepath.Join(homeDir, dir[1:])
	}

	return os.Chdir(dir)
}

func (c *CdCommand) Help() string {
	return `cd: Change the current directory
Usage: cd [directory]

Special paths:
  ..    Move to parent directory
  ~     Move to home directory
  .     Stay in current directory`
}
