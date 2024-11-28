// internal/shell/builtins/ver.go
package builtins

import (
	"fmt"
	"runtime"
)

type VerCommand struct{}

var (
	Version    = "0.1.0"
	CommitHash = "unknown"
	BuildTime  = "unknown"
)

func (v *VerCommand) Execute(args []string) error {
	fmt.Printf("gosh version %s (%s)\n", Version, runtime.Version())
	fmt.Printf("Build: %s\n", BuildTime)
	fmt.Printf("Commit: %s\n", CommitHash)
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	return nil
}

func (v *VerCommand) Help() string {
	return "ver: Display version information for gosh"
}
