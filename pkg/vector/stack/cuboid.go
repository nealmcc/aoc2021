package stack

import (
	"errors"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// Cuboid holds a stack of vector.Cuboid.
// The zero value is empty and ready to use.
type Cuboid struct {
	data []v.Cuboid
}

// Length returns the number of items on the stack.
func (s *Cuboid) Length() int {
	return len(s.data)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *Cuboid) Peek() (v.Cuboid, bool) {
	if len(s.data) > 0 {
		return s.data[len(s.data)-1], true
	}
	return v.Cuboid{}, false
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *Cuboid) Pop() v.Cuboid {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Push adds the given item to the top of the stack
func (s *Cuboid) Push(x v.Cuboid) {
	s.data = append(s.data, x)
}
