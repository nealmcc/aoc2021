package main

import (
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

	// diag, err := render(segments)
	// 	picture := fmt.Sprintf("%v", diag)
	// 	a.Equal(`.......1..
	// ..1....1..
	// ..1....1..
	// .......1..
	// .112111211
	// ..........
	// ..........
	// ..........
	// ..........
	// 222111....`, picture)

	// p1 := part1(diag)
	// a.Equal(5, p1)
}
