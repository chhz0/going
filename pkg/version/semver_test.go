package version

import "testing"

func Test_ParseSemVer(t *testing.T) {
	semver := parseSemVer("v1.2.3-alpha.1+001")
	t.Logf("%#v\n", semver)
}
