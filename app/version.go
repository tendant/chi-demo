package app

import (
	"net/http"
	"runtime/debug"
)

//  As long as you have Go 1.18 or higher, a simple go build should
//  suffice to pass the git information into the Commit string
//  variable.

// Apart from vsc.revision, you can also query for vcs.time (that is
// the timestamp of the commit message) and check vcs.modified (that
// is true if the binary builds from uncommitted changes). To make
// this work, you need Go 1.18, and should build using the -buildvcs
// (which should be available in your goenv).

var (
	GitCommit  = "none"
	CommitTime = ""
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return GitCommit
}()

var Timestamp = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}

	return CommitTime
}()

// for _, kv := range info.Settings {
// 	switch kv.Key {
// 	case "vcs.revision":
// 		Revision = kv.Value
// 	case "vcs.time":
// 		LastCommit, _ = time.Parse(time.RFC3339, kv.Value)
// 	case "vcs.modified":
// 		DirtyBuild = kv.Value == "true"
// 	}
// }

// Middleware that accepts the application version as a parameter
func Version(appVersion string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the X-App-Version header with the provided appVersion
			w.Header().Set("X-App-Version", appVersion)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
