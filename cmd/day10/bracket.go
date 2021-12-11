package main

import (
	"fmt"

	"github.com/nealmcc/aoc2021/pkg/stack"
)

type bracket = byte

type chunk = string

type line struct {
	chunks []chunk
	raw    []bracket
}

var (
	// pairs is a map from the open bracket to its closing bracket
	pairs = map[bracket]bracket{
		'(': ')',
		'[': ']',
		'{': '}',
		'<': '>',
	}
	// errScore is a map from a closing bracket to its value in part 1
	errScore = map[bracket]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}
)

// parse interprets the raw bytes for the line as one or more chunks.
// If the line is invalid, parse will return an error.
func (l *line) parse() (err error) {
	defer func() {
		fmt.Printf("%-110s ", string(l.raw))

		if len(l.chunks) > 0 {
			fmt.Println("found chunks:")
			for _, ch := range l.chunks {
				fmt.Printf("%s\n", string(ch))
			}
		}

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("[ok]")
		}
	}()

	l.chunks = make([]chunk, 0)

	s := &stack.ByteStack{}
	i := 0
	for j, b := range l.raw {
		switch {
		case isOpen(b):
			s.Push(b)

		case isClose(b):
			ok, top := s.Peek()
			if !ok {
				return errCorrupted{pos: j + 1, got: b}
			}
			want := pairs[top]
			if want != b {
				return errCorrupted{pos: j + 1, want: want, got: b}
			}
			s.Pop()
			if s.Length() == 0 {
				chunk := chunk(l.raw[i : j+1])
				l.chunks = append(l.chunks, chunk)
				i = j + 1
			}

		default:
			return fmt.Errorf("invalid input at column %d - got %c", j+1, b)
		}
	}

	if s.Length() > 0 {
		return errIncomplete{pos: len(l.raw) + 1}
	}

	return nil
}

func isOpen(b bracket) bool {
	_, ok := pairs[b]
	return ok
}

func isClose(b bracket) bool {
	_, ok := errScore[b]
	return ok
}

// suggest returns the sequence of brackets that will complete an incomplete line.
// Assumes the line is not corrupt.
func (l *line) suggest() []bracket {
	start := 0
	for _, ch := range l.chunks {
		start += len(ch)
	}

	result := make([]bracket, 0, 8)
	s := &stack.ByteStack{}
	for _, b := range l.raw[start:] {
		switch {
		case isOpen(b):
			s.Push(b)

		case isClose(b):
			top := s.Pop()
			want := pairs[top]
			for want != b {
				result = append(result, want)
				top = s.Pop()
				want = pairs[top]
			}
		}
	}

	for s.Length() > 0 {
		top := s.Pop()
		want := pairs[top]
		result = append(result, want)
	}
	return result
}

func value(suffix []bracket) int {
	value := map[bracket]int{
		')': 1,
		']': 2,
		'}': 3,
		'>': 4,
	}

	score := 0
	for _, b := range suffix {
		score *= 5
		score += value[b]
	}

	return score
}
