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
	if s[i].X < s[j].X {
		return true
	}

	if s[i].X > s[j].X {
		return false
	}

	if s[i].Y < s[j].Y {
		return true
	}

	if s[i].Y > s[j].Y {
		return false
	}

	return s[i].Z < s[j].Z
}
