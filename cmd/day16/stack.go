package main

import (
	"errors"
)

// PacketStack holds a stack of Packets.
// The zero value is empty and ready to use.
type PacketStack struct {
	data []Packet
}

// Length returns the number of items on the stack.
func (s *PacketStack) Length() int {
	return len(s.data)
}

// Peek looks at the top item on the stack, return true and a copy of the item
// if there is one. Returns false if the stack is empty.
func (s *PacketStack) Peek() (bool, Packet) {
	if len(s.data) > 0 {
		return true, s.data[len(s.data)-1]
	}
	return false, Packet{}
}

// Pop removes the top item from the stack, and returns it.
// If the stack is empty, Pop() panics.
func (s *PacketStack) Pop() Packet {
	if len(s.data) == 0 {
		panic(errors.New("cannot pop an empty stack"))
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Push adds the given item to the top of the stack
func (s *PacketStack) Push(x Packet) {
	s.data = append(s.data, x)
}
