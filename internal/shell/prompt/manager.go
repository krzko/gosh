// internal/shell/prompt/manager.go
package prompt

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/chzyer/readline"
	"github.com/krzko/gosh/internal/shell/command"
	"github.com/krzko/gosh/internal/shell/completion"
	"github.com/krzko/gosh/internal/utils/color"
)

type Manager struct {
	format   string
	theme    *color.Theme
	rl       *readline.Instance
	builtins map[string]command.BuiltinCommand
}

type Config struct {
	ShowGitBranch bool
	ShowHostname  bool
	ShowPath      bool
	Theme         string
}

// internal/shell/prompt/manager.go
func NewManager(completer *completion.Completer, builtins map[string]command.BuiltinCommand) (*Manager, error) {
	rlConfig := &readline.Config{
		Prompt:          "> ",
		HistoryFile:     homeDir() + "/.gosh_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:    completer,
	}

	rl, err := readline.NewEx(rlConfig)
	if err != nil {
		return nil, err
	}

	return &Manager{
		format:   "${user}@${hostname}:${pwd}$ ",
		theme:    color.DefaultTheme(),
		rl:       rl,
		builtins: builtins,
	}, nil
}

func (m *Manager) SetFormat(format string) {
	m.format = format
}

func (m *Manager) Read() (string, error) {
	prompt := m.buildPrompt()
	m.rl.SetPrompt(prompt)
	line, err := m.rl.Readline()
	if err != nil {
		return line, err
	}

	// Colorize the command if it's valid
	words := strings.Fields(line)
	if len(words) > 0 {
		cmd := words[0]
		if _, ok := m.builtins[cmd]; ok {
			// Save original command for return
			original := line
			// Show colorized version
			fmt.Printf("\033[1A\033[2K\r%s\033[32m%s\033[0m%s\n",
				prompt, cmd, line[len(cmd):])
			return original, nil
		}
	}

	return line, nil
}

func (m *Manager) buildPrompt() string {
	promptStr := m.format

	// Replace user
	if strings.Contains(promptStr, "${user}") {
		if u, err := user.Current(); err == nil {
			promptStr = strings.ReplaceAll(promptStr, "${user}", u.Username)
		}
	}

	// Replace hostname
	if strings.Contains(promptStr, "${hostname}") {
		if hostname, err := os.Hostname(); err == nil {
			promptStr = strings.ReplaceAll(promptStr, "${hostname}", hostname)
		}
	}

	// Replace pwd
	if strings.Contains(promptStr, "${pwd}") {
		if pwd, err := os.Getwd(); err == nil {
			promptStr = strings.ReplaceAll(promptStr, "${pwd}", pwd)
		}
	}

	return m.theme.ColorizePrompt(promptStr)
}

func (m *Manager) Close() error {
	if m.rl != nil {
		return m.rl.Close()
	}
	return nil
}

func homeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return os.TempDir()
}
