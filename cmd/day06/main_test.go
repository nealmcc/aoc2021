package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = "3,4,3,1,2"

func Test_part1(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	fish, err := read(example)
	r.NoError(err)
	t.Log("fish", fish)
	a.Equal(int64(5), fish.count())
	fish.cycle(80)
	a.Equal(int64(5934), fish.count())
}
