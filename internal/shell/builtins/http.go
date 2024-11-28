// internal/shell/builtins/http.go
package builtins

import (
	"fmt"
)

// HttpCommand represents the 'http' builtin command.
type HttpCommand struct{}

// Execute performs an HTTP GET request to the specified URL.
// Usage: http [URL]
func (c *HttpCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("URL required")
	}

	rawURL := args[0]
	if err := executeRequest(rawURL, "http"); err != nil {
		return err
	}

	return nil
}

// Help returns the help message for the 'http' command.
func (c *HttpCommand) Help() string {
	return `http: Transfer data from a URL using HTTP
Usage: http [URL]

If the scheme (http:// or https://) is omitted, 'http://' is assumed.`
}
