// Package verinfo reports build/version information via the Go stdlib
// (runtime/debug build settings), replacing veg's cmd/veg/verinfo package.
package verinfo

import (
	"runtime"
	"runtime/debug"
)

var (
	Version   = "Local"
	Commit    = "Dirty"
	BuildDate = "Unknown"

	GoVersion = "Unknown"
	BuildOS   = "Unknown"
	BuildArch = "Unknown"

	// the version we can fetch with veg mod; injected into templates
	VegVersion       = "v0.6.10"
	VegSchemaVersion = "v0.0.2"

	CueVersion = "Unknown"
)

func init() {
	info, _ := debug.ReadBuildInfo()
	GoVersion = info.GoVersion
	BuildOS = runtime.GOOS
	BuildArch = runtime.GOARCH

	dirty := false
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			Commit = s.Value
		case "vcs.time":
			BuildDate = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				dirty = true
			}
		}
	}
	if dirty {
		Commit += "+dirty"
	}

	if Version != "Local" {
		VegVersion = Version
	}

	for _, d := range info.Deps {
		if d.Path == "cuelang.org/go" {
			CueVersion = d.Version
		}
	}
}
