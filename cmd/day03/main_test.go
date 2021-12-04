package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	a := assert.New(t)

	m, err := read(strings.NewReader(example))
	require.NoError(t, err)

	a.Equal([]int{7, 5, 8, 7, 5}, m.ones)
	a.Equal(12, m.count)
	a.Equal(int64(22), m.gamma())
	a.Equal(int64(0b_11111), m.maxSample())
	a.Equal(int64(9), m.epsilon())
}

func Test_part2(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	m, err := read(strings.NewReader(example))
	r.NoError(err)

	o2, err := m.oxygen()
	r.NoError(err)
	a.Equal(int64(23), o2)

	co2, err := m.carbonDioxide()
	r.NoError(err)
	a.Equal(int64(10), co2)
}
