package main

import (
	"strings"
	"testing"

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

	a.Equal(40, part1(cave))
}
