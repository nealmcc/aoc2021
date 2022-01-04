package main

import "errors"

// Queue holds a queue of blocks.
// The zero value is empty and ready to use.
type Queue struct {
	data []block
}

// Length returns the number of items in the queue.
func (q *Queue) Length() int { return len(q.data) }

// Peek looks at the item at the front of the queue; returns true and a copy of
// the item if there is one. Returns false if the queue is empty.
func (q *Queue) Peek() (block, bool) {
	if len(q.data) > 0 {
		return q.data[0], true
	}
	return block{}, false
}

// Pop removes the item from the front of the queue and returns it.
// If the queue is empty, Pop() panics.
func (q *Queue) Pop() block {
	if len(q.data) == 0 {
		panic(errors.New("cannot pop an empty queue"))
	}
	val := q.data[0]
	(*q).data = q.data[1:]
	return val
}

// Push adds the given item to the back of the queue
func (q *Queue) Push(x block) {
	q.data = append(q.data, x)
}
