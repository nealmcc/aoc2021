package main

import "errors"

type Stack struct {
	data []state
}

type state struct {
	cave  *node
	rules caveRules
}

// Length returns the number of items on the stack.
func (s *Stack) Length() int {
	return len(s.data)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *Stack) Peek() (bool, state) {
	if len(s.data) > 0 {
		return true, s.data[len(s.data)-1]
	}
	return false, state{}
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *Stack) Pop() state {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Push adds the given item to the top of the stack
func (s *Stack) Push(x state) {
	s.data = append(s.data, x)
}
