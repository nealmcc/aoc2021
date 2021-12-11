// Package stack implements a stack of bytes.
package stack

import "errors"

// Stack holds a stack of bytes. The zero value is empty and ready to use.
type Stack struct {
	data []byte
}

// Push adds the given byte to the top of the stack
func (s *Stack) Push(b byte) {
	s.data = append(s.data, b)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *Stack) Peek() (bool, byte) {
	if len(s.data) > 0 {
		return true, s.data[len(s.data)-1]
	}
	return false, 0
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *Stack) Pop() byte {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Length returns the number of items on the stack.
func (s *Stack) Length() int {
	return len(s.data)
}
