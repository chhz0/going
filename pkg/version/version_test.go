package version

import (
	"encoding/json"
	"testing"
)

func TestString(t *testing.T) {
	version = "v1.2.3"
	if got := String(); got != "v1.2.3" {
		t.Errorf("String() = %v, want %v", got, "v1.2.3")
	}
}

func TestJSON(t *testing.T) {
	version = "v1.2.3"
	gitCommit = "abc123"
	buildDate = "2023-01-01T12:00:00Z"

	jsonStr, err := JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}

	var info Info
	if err := json.Unmarshal([]byte(jsonStr), &info); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if info.Version != "v1.2.3" {
		t.Errorf("JSON Version = %v, want v1.2.3", info.Version)
	}
	if info.GitCommit != "abc123" {
		t.Errorf("JSON GitCommit = %v, want abc123", info.GitCommit)
	}
}

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		input    string
		expected *SemVer
	}{
		{"v1.2.3", &SemVer{Major: "1", Minor: "2", Patch: "3", Original: "1.2.3"}},
		{"v1.2.3-alpha.1", &SemVer{Major: "1", Minor: "2", Patch: "3", Prerelease: "alpha.1", Original: "1.2.3-alpha.1"}},
		{"v1.2.3+build.123", &SemVer{Major: "1", Minor: "2", Patch: "3", Build: "build.123", Original: "1.2.3+build.123"}},
		{"v1.2.3-alpha.1+build.123", &SemVer{Major: "1", Minor: "2", Patch: "3", Prerelease: "alpha.1", Build: "build.123", Original: "1.2.3-alpha.1+build.123"}},
		{"invalid", nil},
	}

	for _, test := range tests {
		got := parseSemVer(test.input)
		if test.expected == nil {
			if got != nil {
				t.Errorf("parseSemVer(%v) = %v, want nil", test.input, got)
			}
			continue
		}

		if got == nil {
			t.Errorf("parseSemVer(%v) = nil, want %v", test.input, test.expected)
			continue
		}

		if got.Major != test.expected.Major || got.Minor != test.expected.Minor || got.Patch != test.expected.Patch {
			t.Errorf("parseSemVer(%v) = %v, want %v", test.input, got, test.expected)
		}
	}
}

// func TestIsDevVersion(t *testing.T) {
// 	tests := []struct {
// 		version string
// 		want    bool
// 	}{
// 		{"development", true},
// 		{"v1.0.0-dev", true},
// 		{"v1.0.0-alpha.1", true},
// 		{"v1.0.0-beta.2", true},
// 		{"v1.0.0-rc.1", true},
// 		{"v1.0.0", false},
// 	}

// 	for _, test := range tests {
// 		version = test.version
// 		if got := IsDevVersion(); got != test.want {
// 			t.Errorf("IsDevVersion(%v) = %v, want %v", test.version, got, test.want)
// 		}
// 	}
// }

// func TestCompare(t *testing.T) {
// 	tests := []struct {
// 		v1, v2 string
// 		want   int
// 	}{
// 		{"v1.0.0", "v1.0.0", 0},
// 		{"v1.0.1", "v1.0.0", 1},
// 		{"v1.0.0", "v1.0.1", -1},
// 		{"v1.1.0", "v1.0.9", 1},
// 		{"v2.0.0", "v1.9.9", 1},
// 	}

// 	for _, test := range tests {
// 		if got := Compare(test.v1, test.v2); got != test.want {
// 			t.Errorf("Compare(%v, %v) = %v, want %v", test.v1, test.v2, got, test.want)
// 		}
// 	}
// }
