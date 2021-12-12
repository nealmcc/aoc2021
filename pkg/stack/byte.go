package stack

import "errors"

// ByteStack holds a stack of bytes. The zero value is empty and ready to use.
type ByteStack struct {
	data []byte
}

// Length returns the number of items on the stack.
func (s *ByteStack) Length() int {
	return len(s.data)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *ByteStack) Peek() (bool, byte) {
	if len(s.data) > 0 {
		return true, s.data[len(s.data)-1]
	}
	return false, 0
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *ByteStack) Pop() byte {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Push adds the given item to the top of the stack
func (s *ByteStack) Push(x byte) {
	s.data = append(s.data, x)
}
