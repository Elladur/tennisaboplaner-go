package internal

// Version is the current version of tennisaboplaner-go.
// It is overridden at build time via:
//
//	-ldflags "-X github.com/Elladur/tennisaboplaner-go/internal.Version=vX.Y.Z"
var Version = "dev"

const (
	// RepoOwner is the GitHub owner of the project.
	RepoOwner = "Elladur"
	// RepoName is the GitHub repository name.
	RepoName = "tennisaboplaner-go"
)
