package version

import (
	"fmt"
	"os"
	"strings"
)

var (
	// GoVersion go version, setup by makefile
	GoVersion = ""
	// BranchName git branch, setup by makefile
	BranchName = ""
	// CommitID git commit id, setup by makefile
	CommitID = ""
	// BuildTime build time, setup by makefile
	BuildTime = ""
	// Version version, setup by makefile
	Version = ""
)

func GenVersionMessage() string {
	messages := []string{
		"Go version:\t" + GoVersion,
		"Branch name:\t" + BranchName,
		"Git commit:\t" + CommitID,
		"Build time:\t" + BuildTime,
		"Version:\t" + Version,
	}
	return strings.Join(messages, "\n")
}

func Print() {
	fmt.Fprintln(os.Stdout, GenVersionMessage())
}
