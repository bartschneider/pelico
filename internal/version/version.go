package version

import (
	"fmt"
	"time"
)

const (
	Major = 1
	Minor = 2
	Patch = 2
)

var (
	// BuildTime will be set during build
	BuildTime = ""
	// GitCommit will be set during build  
	GitCommit = ""
)

type Info struct {
	Version   string    `json:"version"`
	BuildTime string    `json:"build_time"`
	GitCommit string    `json:"git_commit"`
	Timestamp time.Time `json:"timestamp"`
}

func GetVersion() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}

func GetInfo() Info {
	return Info{
		Version:   GetVersion(),
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		Timestamp: time.Now(),
	}
}
