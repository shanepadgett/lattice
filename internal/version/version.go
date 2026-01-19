package version

import (
	"runtime/debug"
	"strconv"
	"strings"
)

// ParseSemverMajor returns the major version for tags like "v1.2.3".
func ParseSemverMajor(v string) (int, bool) {
	v = strings.TrimSpace(v)
	if !strings.HasPrefix(v, "v") {
		return 0, false
	}

	// Parse digits until '.' (or end).
	start := 1
	end := start
	for end < len(v) {
		c := v[end]
		if c == '.' {
			break
		}
		if c < '0' || c > '9' {
			return 0, false
		}
		end++
	}
	if end == start {
		return 0, false
	}

	major, err := strconv.Atoi(v[start:end])
	if err != nil || major <= 0 {
		return 0, false
	}
	return major, true
}

// BinaryMajorVersion tries to determine the major version of the current binary.
//
// Prefer tags embedded by Go's VCS info (setting "vcs.tag"). If unavailable,
// fall back to the provided mainVersion (typically set via -ldflags).
func BinaryMajorVersion(mainVersion string) (int, bool) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key != "vcs.tag" {
				continue
			}
			if major, ok := ParseSemverMajor(setting.Value); ok {
				return major, true
			}
		}
	}

	return ParseSemverMajor(mainVersion)
}
