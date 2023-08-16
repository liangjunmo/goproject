package version

import "fmt"

var (
	BuildDate string
	GoVersion string
	GitCommit string
)

const Version = "dev"

func Describe() string {
	return fmt.Sprintf(`
Version: %s
BuildDate: %s
GoVersion: %s
GitCommit: %s
`, Version, BuildDate, GoVersion, GitCommit)
}
