package cmd

import (
	"runtime/debug"
)

var version = "unknown"

func getVersion() string {
	if version == "unknown" {
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			version = buildInfo.Main.Version
		}
	}

	return version
}
