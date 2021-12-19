package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	nodes, err := read(strings.NewReader("[[1,2],[[3,4],5]]"))
	r.NoError(err)

	r.Equal(1, len(nodes))
	a.Equal(15, nodes[0].Value())
	a.Equal(143, nodes[0].Magnitude())
}
