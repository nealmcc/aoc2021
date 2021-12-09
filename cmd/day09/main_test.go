package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var example string = `2199943210
3987894921
9856789892
8767896789
9899965678
`

func Test_read(t *testing.T) {
	r := require.New(t)

	got, err := read(strings.NewReader(example))
	r.NoError(err)

	r.Equal(50, len(got))
}

func Test_part1(t *testing.T) {
	r := require.New(t)

	land, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part1(land)

	r.Equal(15, got)
}

func Test_part2(t *testing.T) {
	r := require.New(t)

	land, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part2(land)

	r.Equal(1134, got)
}
