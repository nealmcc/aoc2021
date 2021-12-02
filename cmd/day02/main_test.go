package main

import (
	"strings"
	"testing"

	"github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/require"
)

var example string = `forward 5
down 5
forward 8
up 3
down 8
forward 2
`

func Test_part1(t *testing.T) {
	r := require.New(t)

	got, err := part1(vector.Coord{}, strings.NewReader(example))
	r.NoError(err)
	r.Equal(vector.Coord{X: 15, Y: 10}, got)
}
