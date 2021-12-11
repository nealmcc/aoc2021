package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `5483143223
2745854711
5264556173
6141336146
6357385478
4167524645
2176841721
6882881134
4846848554
5283751526
`

func Test_read(t *testing.T) {
	r := require.New(t)

	grid, err := read(strings.NewReader(example))
	r.NoError(err)

	r.Equal(100, len(grid))
}

func Test_step(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	grid, err := read(strings.NewReader(`11111
19991
19191
19991
11111`))
	r.NoError(err)

	n := grid.step()
	a.Equal(9, n)

	want, err := read(strings.NewReader(`34543
40004
50005
40004
34543`))
	r.NoError(err)

	r.Equal(want, grid)
}

func Test_part1(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	grid, err := read(strings.NewReader(example))
	r.NoError(err)

	n := part1(&grid)
	a.Equal(1656, n)

	want, err := read(strings.NewReader(`0397666866
0749766918
0053976933
0004297822
0004229892
0053222877
0532222966
9322228966
7922286866
6789998766`))
	r.NoError(err)

	a.Equal(want, grid)
}

func Test_part2(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	grid, err := read(strings.NewReader(example))
	r.NoError(err)

	n := part2(&grid)
	a.Equal(195, n)
}
