// internal/shell/parser/parser.go
package parser

import (
	"errors"
	"strings"

	"github.com/krzko/gosh/internal/shell/executor"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) (*executor.Command, error) {
	if input == "" {
		return nil, errors.New("empty input")
	}

	// Handle ".." as "cd .."
	if input == ".." {
		return &executor.Command{
			Name: "cd",
			Args: []string{".."},
		}, nil
	}

	// Split the input into pipeline segments
	pipeSegments := strings.Split(input, "|")
	var currentCommand *executor.Command
	var firstCommand *executor.Command

	for i, segment := range pipeSegments {
		// Trim spaces and split into command and args
		parts := splitCommand(strings.TrimSpace(segment))
		if len(parts) == 0 {
			continue
		}

		command := &executor.Command{
			Name: parts[0],
			Args: parts[1:],
		}

		if i == 0 {
			firstCommand = command
			currentCommand = command
		} else {
			currentCommand.Pipe = command
			currentCommand = command
		}
	}

	return firstCommand, nil
}

// splitCommand splits a command string into command and arguments,
// handling quoted strings and escaping
func splitCommand(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	escapeNext := false

	for _, char := range input {
		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\\':
			escapeNext = true
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if !inQuotes {
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
