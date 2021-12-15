package main

import (
	"strings"
	"testing"

	"github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581
`

func Test_part1(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	cave, err := read(strings.NewReader(example))
	r.NoError(err)

	from := vector.Coord{X: 0, Y: 0}
	to := vector.Coord{X: 9, Y: 9}

	a.Equal(40, shortestPath(cave, from, to))
}
