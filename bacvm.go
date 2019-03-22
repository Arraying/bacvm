package bacvm

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// VersionMajor is the major version.
	VersionMajor = 1
	// VersionMinor is the minor version.
	VersionMinor = 0
)

// Offer offers the VM a version, which can be accepted or rejected.
func Offer(version string) bool {
	split := strings.SplitN(version, ".", 2)
	if len(split) != 2 {
		return false
	}
	major, err := strconv.Atoi(split[0])
	if err != nil || VersionMajor != major {
		return false
	}
	_, err = strconv.Atoi(split[0])
	if err != nil {
		return false
	}
	return true
}

// Version gets the version as a string.
func Version() string {
	return fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
}
