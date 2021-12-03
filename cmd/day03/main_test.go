package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var example string = `00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010
`

func Test_part1(t *testing.T) {
	r := require.New(t)

	m, err := part1(strings.NewReader(example))
	r.NoError(err)
	r.Equal(&meter{
		ones:  []int{7, 5, 8, 7, 5},
		count: 12,
	}, m)
	r.Equal(int64(22), m.gamma())
	r.Equal(int64(0b_11111), m.maxSample())
	r.Equal(int64(9), m.epsilon())
}
