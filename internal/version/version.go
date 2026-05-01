package version

// Version is injected at build time via ldflags:
//
//	-X github.com/pantheon-org/iris/internal/version.Version=<tag>
var Version = "dev"
