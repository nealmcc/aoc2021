package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var example string = `199
200
208
210
200
207
240
269
260
263
`

func Test_part1(t *testing.T) {
	r := require.New(t)

	got, err := part1(strings.NewReader(example))
	r.NoError(err)
	r.Equal(7, got)
}
