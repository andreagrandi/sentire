package version

import (
	"fmt"
	"runtime"
)

// Version is the current version of sentire
const Version = "0.1.0"

// These variables will be set at build time via ldflags
var (
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// GetVersionInfo returns formatted version information
func GetVersionInfo() string {
	return fmt.Sprintf("sentire version %s", Version)
}

// GetFullVersionInfo returns detailed version information
func GetFullVersionInfo() string {
	return fmt.Sprintf(`sentire version %s
Build time: %s
Git commit: %s
Go version: %s
OS/Arch: %s/%s`,
		Version,
		BuildTime,
		GitCommit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
