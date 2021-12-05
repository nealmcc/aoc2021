package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7
`

func Test_part1(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	turns, boards, err := read(strings.NewReader(example))
	r.NoError(err)

	a.Equal([]int{
		7, 4, 9, 5, 11, 17, 23, 2, 0, 14, 21, 24, 10, 16,
		13, 6, 15, 25, 12, 22, 18, 20, 8, 19, 3, 26, 1,
	}, turns)

	a.Equal(`22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19`, fmt.Sprintf("%v", boards[0]))

	i, n := part1(boards, turns)
	a.Equal(2, i)
	a.Equal(11, n)
	a.Equal(188, boards[i].sum())
}
