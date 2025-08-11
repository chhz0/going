package version

import "strings"

type SemVer struct {
	Major      string
	Minor      string
	Patch      string
	Prerelease string
	Build      string
	Original   string
}

func parseSemVer(version string) *SemVer {
	version = strings.TrimPrefix(version, "v")

	parts := strings.SplitN(version, "+", 2)
	mainVersin := parts[0]
	build := ""
	if len(parts) > 1 {
		build = parts[1]
	}

	parts = strings.SplitN(mainVersin, "-", 2)
	versionCore := parts[0]
	prerelease := ""
	if len(parts) > 1 {
		prerelease = parts[1]
	}

	components := strings.Split(versionCore, ".")
	if len(components) != 3 {
		return nil
	}

	return &SemVer{
		Major:      components[0],
		Minor:      components[1],
		Patch:      components[2],
		Prerelease: prerelease,
		Build:      build,
		Original:   version,
	}
}

func IsDevVersion() bool {
	return !strings.Contains(gitCommit, "dev")
}
