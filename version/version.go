package version

import (
	"fmt"
)

const (
	Major = 0
	Minor = 1
	Patch = 0
)

// String formats the version string for this module in semver format.
func String() string {
	return fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
}
