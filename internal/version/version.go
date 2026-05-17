// Package version provides runtime version comparison utilities,
// allowing services to advertise and validate their minimum required
// grpc-healthd version.
package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Version represents a semantic version with major, minor, and patch components.
type Version struct {
	Major int
	Minor int
	Patch int
}

// Parse parses a semantic version string of the form "vMAJOR.MINOR.PATCH" or
// "MAJOR.MINOR.PATCH". It returns an error if the string is malformed.
func Parse(s string) (Version, error) {
	s = strings.TrimPrefix(s, "v")
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("version: invalid format %q", s)
	}
	var v Version
	var err error
	v.Major, err = strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, errors.New("version: invalid major component")
	}
	v.Minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, errors.New("version: invalid minor component")
	}
	v.Patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, errors.New("version: invalid patch component")
	}
	return v, nil
}

// String returns the canonical string representation, e.g. "1.2.3".
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// AtLeast reports whether v is greater than or equal to other.
func (v Version) AtLeast(other Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch >= other.Patch
}

// Compare returns -1, 0, or 1 if v is less than, equal to, or greater than other.
func (v Version) Compare(other Version) int {
	switch {
	case v.AtLeast(other) && other.AtLeast(v):
		return 0
	case v.AtLeast(other):
		return 1
	default:
		return -1
	}
}
