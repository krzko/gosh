// internal/shell/builtins/https.go
package builtins

import (
	"fmt"
)

// HttpsCommand represents the 'https' builtin command.
type HttpsCommand struct{}

// Execute performs an HTTPS GET request to the specified URL.
// Usage: https [URL]
func (c *HttpsCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("URL required")
	}

	rawURL := args[0]
	if err := executeRequest(rawURL, "https"); err != nil {
		return err
	}

	return nil
}

// Help returns the help message for the 'https' command.
func (c *HttpsCommand) Help() string {
	return `https: Transfer data from a URL using HTTPS
Usage: https [URL]

If the scheme (https://) is omitted, 'https://' is assumed.`
}
