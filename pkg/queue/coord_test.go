package queue

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoord_pushAndLen(t *testing.T) {
	q := make(Coord, 0)

	nodes := []*CoordNode{
		{
			Value:    vector.Coord{X: 4, Y: 2},
			Priority: 42,
		},
		{
			Value:    vector.Coord{X: 6, Y: 7},
			Priority: 67,
		},
		{
			Value:    vector.Coord{X: 0, Y: 0},
			Priority: 0,
		},
	}

	for _, n := range nodes {
		q.Push(n)
	}

	require.Equal(t, 3, q.Len())
}

func TestCoord_initAndPop(t *testing.T) {
	nodes := []*CoordNode{
		{
			Value:    vector.Coord{X: 4, Y: 2},
			Priority: 42,
		},
		{
			Value:    vector.Coord{X: 6, Y: 7},
			Priority: 67,
		},
		{
			Value:    vector.Coord{X: 0, Y: 0},
			Priority: 0,
		},
	}

	q := Coord(nodes)
	fmt.Println(q)

	heap.Init(&q)
	fmt.Println(q)

	node, ok := heap.Pop(&q).(*CoordNode)
	require.True(t, ok)
	assert.Equal(t, vector.Coord{X: 6, Y: 7}, node.Value)
	assert.Equal(t, 67, node.Priority)

	node, ok = heap.Pop(&q).(*CoordNode)
	require.True(t, ok)
	assert.Equal(t, vector.Coord{X: 4, Y: 2}, node.Value)
	assert.Equal(t, 42, node.Priority)
}

func TestCoord_update(t *testing.T) {
	items := map[vector.Coord]int{
		{X: 4, Y: 2}: 42,
		{X: 6, Y: 7}: 67,
		{X: 0, Y: 0}: 0,
	}

	q := new(Coord)
	for val, prio := range items {
		heap.Push(q, &CoordNode{Value: val, Priority: prio})
	}

	newItem := &CoordNode{Value: vector.Coord{X: 99, Y: 99}}
	// push the item on to the queue (it will have priority 0)
	heap.Push(q, newItem)

	// now update the item's priority:
	q.Update(newItem, newItem.Value, 9999)

	node, ok := heap.Pop(q).(*CoordNode)
	require.True(t, ok)
	assert.Equal(t, vector.Coord{X: 99, Y: 99}, node.Value)
	assert.Equal(t, 9999, node.Priority)
}
