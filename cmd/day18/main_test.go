package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	tokens, err := read(strings.NewReader("[[1,2],[[3,4],5]]"))
	r.NoError(err)

	a.Equal(int64(17), tokens.Size())

	// a.Equal(143, tree.Magnitude())
}
