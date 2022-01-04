package main

import (
	"testing"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueuePushPeekPopLength(t *testing.T) {
	r, a := require.New(t), assert.New(t)
	q := Queue{}

	b0 := block{
		Sensors: map[int]v.I3{0: {}},
		Beacons: []v.I3{},
	}

	b1 := block{
		Sensors: map[int]v.I3{1: {X: 1, Y: 1}},
		Beacons: []v.I3{{X: 1, Y: 1}},
	}

	b2 := block{
		Sensors: map[int]v.I3{2: {X: 2, Y: 2}},
		Beacons: []v.I3{
			{X: 2, Y: 2},
		},
	}

	q.Push(b0)
	q.Push(b1)
	q.Push(b2)

	r.Equal(3, q.Length())

	first, ok := q.Peek()
	r.True(ok)
	a.Equal(3, q.Length())
	a.Equal(b0, first)

	first = q.Pop()
	a.Equal(2, q.Length())
	a.Equal(b0, first)

	second := q.Pop()
	a.Equal(1, q.Length())
	a.Equal(b1, second)

	third := q.Pop()
	a.Equal(0, q.Length())
	a.Equal(b2, third)
}
