package internal

import "runtime/debug"
import "strings"

// Version is the current version of tennisaboplaner-go.
// It is overridden at build time via:
//
//	-ldflags "-X github.com/Elladur/tennisaboplaner-go/internal.Version=vX.Y.Z"
//
// When not overridden (e.g. `go install ...@version`), it falls back to
// the module version embedded by the Go toolchain.
var Version = "dev"

func init() {
	// Normalize an optional leading "v" so ldflags and module versions agree.
	if Version != "dev" {
		Version = strings.TrimPrefix(Version, "v")
		return
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		v := strings.TrimPrefix(info.Main.Version, "v")
		if v != "" && v != "(devel)" {
			Version = v
		}
	}
}

const (
	// RepoOwner is the GitHub owner of the project.
	RepoOwner = "Elladur"
	// RepoName is the GitHub repository name.
	RepoName = "tennisaboplaner-go"
)
