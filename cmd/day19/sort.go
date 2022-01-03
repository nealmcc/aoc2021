package main

import (
	"sort"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// byXYZ implements sort.Interface for []v.I3 based on
// the x coordinate, then y, and then z.
type byXYZ []v.I3

var _ sort.Interface = byXYZ{}

func (s byXYZ) Len() int      { return len(s) }
func (s byXYZ) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byXYZ) Less(i, j int) bool {
	return isLess(s[i], s[j])
}

// isLess returns true iff a < b based on the X, then Y, then Z coordinates.
func isLess(a, b v.I3) bool {
	if a.X < b.X {
		return true
	}

	if a.X > b.X {
		return false
	}

	if a.Y < b.Y {
		return true
	}

	if a.Y > b.Y {
		return false
	}

	return a.Z < b.Z
}
