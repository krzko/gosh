// internal/shell/completion/completer.go
package completion

import (
	"os"
	"path/filepath"
	"strings"
)

type Completer struct {
	builtins []string
}

func NewCompleter() *Completer {
	return &Completer{
		builtins: []string{"cd", "ls", "pwd", "http", "https", "history", "exit", "help"},
	}
}

// Do implements the readline.AutoCompleter interface
func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line[:pos])
	if lineStr == "" {
		return c.completeCommands(""), 0
	}

	words := strings.Fields(lineStr)
	if len(words) == 0 {
		return c.completeCommands(""), 0
	}

	wordStart := 0
	if pos > 0 && !strings.HasSuffix(lineStr, " ") {
		wordStart = strings.LastIndex(lineStr[:pos], " ") + 1
	}

	lastWord := lineStr[wordStart:pos]

	// Command completion
	if len(words) == 1 && lineStr == lastWord {
		return c.completeCommands(lastWord), len(lastWord)
	}

	// Path completion for cd and ls
	if len(words) > 0 && (words[0] == "cd" || words[0] == "ls") {
		if len(words) == 1 || (len(words) == 2 && !strings.HasSuffix(lineStr, " ")) {
			return c.completePaths(lastWord), len(lastWord)
		}
	}

	return nil, 0
}

func (c *Completer) completeCommands(prefix string) [][]rune {
	var suggestions [][]rune
	for _, cmd := range c.builtins {
		if strings.HasPrefix(cmd, prefix) {
			suggestions = append(suggestions, []rune(cmd))
		}
	}
	return suggestions
}

func (c *Completer) completePaths(prefix string) [][]rune {
	var suggestions [][]rune

	dir := "."
	if prefix != "" {
		dir = filepath.Dir(prefix)
		if dir == "." && !strings.HasPrefix(prefix, ".") {
			dir = prefix
		}
	}

	base := filepath.Base(prefix)
	if dir != "." && dir != prefix {
		base = filepath.Base(prefix)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, base) {
			fullPath := name
			if dir != "." {
				fullPath = filepath.Join(dir, name)
			}
			if entry.IsDir() {
				fullPath += "/"
			}
			suggestions = append(suggestions, []rune(fullPath))
		}
	}

	return suggestions
}
