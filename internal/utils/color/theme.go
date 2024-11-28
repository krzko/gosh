// internal/utils/color/theme.go
package color

import (
	"os"

	"github.com/fatih/color"
)

type Theme struct {
	dirColor     *color.Color
	fileColor    *color.Color
	execColor    *color.Color
	symlinkColor *color.Color
	promptColor  *color.Color
}

func DefaultTheme() *Theme {
	return &Theme{
		dirColor:     color.New(color.FgBlue, color.Bold),
		fileColor:    color.New(color.FgWhite),
		execColor:    color.New(color.FgGreen, color.Bold),
		symlinkColor: color.New(color.FgCyan),
		promptColor:  color.New(color.FgYellow, color.Bold),
	}
}

func (t *Theme) ColorizeName(name string, isDir bool, mode os.FileMode) string {
	if isDir {
		return t.dirColor.Sprint(name)
	}
	// Check if file is executable for user, group, or others
	if mode&0111 != 0 {
		return t.execColor.Sprint(name)
	}
	return t.fileColor.Sprint(name)
}

func (t *Theme) ColorizePermissions(perms string) string {
	return t.fileColor.Sprint(perms)
}

func (t *Theme) ColorizePrompt(prompt string) string {
	return t.promptColor.Sprint(prompt)
}
