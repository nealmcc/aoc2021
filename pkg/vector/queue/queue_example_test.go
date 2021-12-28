package queue

import (
	"container/heap"
	"fmt"

	"github.com/nealmcc/aoc2021/pkg/vector"
)

func ExampleCoord_Update() {
	items := map[vector.Coord]int{
		{X: 4, Y: 2}: 42,
		{X: 6, Y: 7}: 67,
		{X: 0, Y: 0}: 0,
	}

	q := new(Coord)
	for val, prio := range items {
		heap.Push(q, &CoordNode{Value: val, Priority: prio})
	}

	// push the item on to the queue (it will have priority 0)
	newItem := &CoordNode{Value: vector.Coord{X: 99, Y: 99}}
	heap.Push(q, newItem)

	// now update the item's priority:
	q.Update(newItem, newItem.Value, 9999)

	node, _ := heap.Pop(q).(*CoordNode)
	fmt.Printf("priority: %2d value: %+v\n", node.Priority, node.Value)
	// Output: priority: 9999 value: {X:99 Y:99}
}
