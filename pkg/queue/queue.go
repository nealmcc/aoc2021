// Package queue implements priority queues.
//
// The zero value is ready to use:
//    q := new(Coord)
//
// Use heap.Push() and heap.Pop() to push and pop items on to the queue:
//    heap.Push(q, &CoordNode{ Value: vector.Coord{X: 1, Y: 2}, Priority: 3})
//    node, ok := heap.Pop(q).(*CoordNode)
//
// Use Update() to update the value and/or priority of an item in the queue:
//    heap.Push(q, node)
//    q.Update(node, node.Value, 4) // keeps the existing value
//
// adapted from the example at: https://pkg.go.dev/container/heap
package queue
