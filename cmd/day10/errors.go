package main

import "fmt"

type ErrIncomplete interface {
	error
	IsIncomplete() bool
}

type ErrCorrupted interface {
	error
	IsCorrupt() bool
	Score() int
}

// compile-time interface checks
var (
	_ ErrIncomplete = errIncomplete{}
	_ ErrCorrupted  = errCorrupted{}
)

type errIncomplete struct {
	pos int
}

func (e errIncomplete) Error() string {
	return fmt.Sprintf("incomplete chunk at position %d", e.pos)
}

func (e errIncomplete) IsIncomplete() bool {
	return true
}

type errCorrupted struct {
	pos       int
	want, got bracket
}

func (e errCorrupted) Error() string {
	if e.want == 0 {
		return fmt.Sprintf("corrupt chunk at position %d - want open ; got %c",
			e.pos, e.got)
	}
	return fmt.Sprintf("corrupt chunk at position %d - want %c ; got %c",
		e.pos, e.want, e.got)
}

func (e errCorrupted) IsCorrupt() bool {
	return true
}

func (e errCorrupted) Score() int {
	return errScore[e.got]
}
