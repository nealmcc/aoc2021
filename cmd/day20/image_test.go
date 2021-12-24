package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	got := fmt.Sprintf("%v", _testExample.img)

	assert.Equal(t, `#..#.
#....
##..#
..#..
..###
`, got)
}
