package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/gosuri/uitable"
)

var (
	version        = ""
	gitCommit      = ""
	gitCommitStamp = ""
	gitBranch      = ""
	gitState       = ""
	buildDate      = ""
)

type Info struct {
	Version       string `json:"version"`
	GitCommit     string `json:"git_commit"`
	GitCommitDate string `json:"git_commit_date,omitempty"`
	GitBranch     string `json:"git_branch"`
	GitState      string `json:"git_state,omitempty"`
	BuildDate     string `json:"build_date"`
	GoVersion     string `json:"go_version"`
	Compiler      string `json:"compiler"`
	Platform      string `json:"platform"`
	Prerelease    string `json:"prerelease,omitempty"`
	BuildMetadata string `json:"build_metadata,omitempty"`
}

func Get() Info {
	info := Info{
		Version:   version,
		GitCommit: gitCommit,
		GitBranch: gitBranch,
		GitState:  gitState,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	if gitCommitStamp != "" {
		if stamp, err := strconv.ParseInt(gitCommitStamp, 10, 64); err == nil {
			info.GitCommitDate = time.Unix(stamp, 0).Format("2006-01-02 15:04:05")
		}
	}

	if v := parseSemVer(version); v != nil {
		info.Prerelease = v.Prerelease
		info.BuildMetadata = v.Build
	}

	return info
}

func String() string {
	return version
}

func Short() string {
	if gitCommit != "" {
		return fmt.Sprintf("%s-%s", version, gitCommit[:7])
	}
	return version
}

func Text() string {
	info := Get()

	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("version", info.Version)
	table.AddRow("git commit", info.GitCommit)
	table.AddRow("git commit data", info.GitCommitDate)
	table.AddRow("git branch", info.GitBranch)
	table.AddRow("git state", info.GitState)
	table.AddRow("build date", info.BuildDate)
	table.AddRow("go version", info.GoVersion)
	table.AddRow("compiler", info.Compiler)
	table.AddRow("platform", info.Platform)

	return table.String()
}

func JSON() (string, error) {
	info := Get()
	data, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
