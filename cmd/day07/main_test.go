package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = "16,1,2,0,4,2,7,1,2,14"

func Test_part1(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	crabs, err := read(example)
	r.NoError(err)
	a.Equal(10, len(crabs))

	pos, fuel := part1(crabs)
	a.Equal(2, pos)
	a.Equal(37, fuel)
}
