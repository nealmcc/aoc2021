package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `6,10
0,14
9,10
0,3
10,4
4,11
6,0
6,12
4,1
0,13
10,12
3,4
3,0
8,4
1,10
2,14
8,10
9,0

fold along y=7
fold along x=5
`

func Test_read(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	dots, folds, err := read(strings.NewReader(example))
	r.NoError(err)

	a.Equal(18, len(dots))
	a.Equal(2, len(folds))
}

func Test_part1(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	dots, folds, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part1(dots, folds[0])
	a.Equal(17, got)
}
