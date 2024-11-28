// internal/shell/builtins/http_utils.go
package builtins

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// executeRequest performs an HTTP GET request to the specified URL and writes the response to stdout.
func executeRequest(rawURL string, defaultScheme string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// If scheme is missing, prepend the default scheme
	if parsedURL.Scheme == "" {
		parsedURL, err = url.Parse(defaultScheme + "://" + rawURL)
		if err != nil {
			return fmt.Errorf("invalid URL after adding default scheme: %w", err)
		}
	}

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return fmt.Errorf("failed to perform GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}

	return nil
}
