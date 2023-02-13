package app

import "runtime/debug"

//  As long as you have Go 1.18 or higher, a simple go build should
//  suffice to pass the git information into the Commit string
//  variable.

// Apart from vsc.revision, you can also query for vcs.time (that is
// the timestamp of the commit message) and check vcs.modified (that
// is true if the binary builds from uncommitted changes). To make
// this work, you need Go 1.18, and should build using the -buildvcs
// (which should be available in your goenv).

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return "NoCommit"
}()

var Timestamp = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}

	return "NoTimestamp"
}()
