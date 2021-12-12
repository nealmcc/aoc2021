package stack

import (
	"errors"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// CoordStack holds a stack of vector.Coord.
// The zero value is empty and ready to use.
type CoordStack struct {
	data []v.Coord
}

// Length returns the number of items on the stack.
func (s *CoordStack) Length() int {
	return len(s.data)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *CoordStack) Peek() (bool, v.Coord) {
	if len(s.data) > 0 {
		return true, s.data[len(s.data)-1]
	}
	return false, v.Coord{}
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *CoordStack) Pop() v.Coord {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Push adds the given item to the top of the stack
func (s *CoordStack) Push(x v.Coord) {
	s.data = append(s.data, x)
}
