package ast

import (
	"errors"
)

// stack holds a stack of Numbers.
// The zero value is empty and ready to use.
type stack struct {
	data []Number
}

// Length returns the number of items on the stack.
func (s *stack) Length() int {
	return len(s.data)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *stack) Peek() (Number, bool) {
	if len(s.data) > 0 {
		return s.data[len(s.data)-1], true
	}
	return nil, false
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *stack) Pop() Number {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Push adds the given item to the top of the stack
func (s *stack) Push(x Number) {
	s.data = append(s.data, x)
}
