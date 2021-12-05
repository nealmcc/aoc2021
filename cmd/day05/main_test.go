package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2
`

func Test_part1(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	segments, err := read(strings.NewReader(example))
	r.NoError(err)
	a.Equal(10, len(segments))

	diagram := render(segments, false)
	a.Equal(`.......1..
..1....1..
..1....1..
.......1..
.112111211
..........
..........
..........
..........
222111....`, fmt.Sprintf("%v", diagram))

	p1 := count(diagram)
	a.Equal(5, p1)
}

func Test_part2(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	segments, err := read(strings.NewReader(example))
	r.NoError(err)
	a.Equal(10, len(segments))

	diagram := render(segments, true)
	a.Equal(`1.1....11.
.111...2..
..2.1.111.
...1.2.2..
.112313211
...1.2....
..1...1...
.1.....1..
1.......1.
222111....`, fmt.Sprintf("%v", diagram))

	p2 := count(diagram)
	a.Equal(12, p2)
}
