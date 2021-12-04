package main

import (
	"fmt"
	"strings"
)

type board struct {
	values  map[int]*square
	stamped [5][5]bool
	won     bool
}

// compile-time interface check
var _ fmt.Formatter = new(board)

// Format implements fmt.Formatter. It ignores the formatting state and verb.
func (b *board) Format(s fmt.State, verb rune) {
	symbols := [5][5]string{}
	for k, v := range b.values {
		s := fmt.Sprintf("%2d", k)
		symbols[v.row][v.col] = s
	}
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if b.stamped[r][c] {
				symbols[r][c] = " x"
			}
		}
	}
	for i := 0; i < 4; i++ {
		s.Write([]byte(strings.Join(symbols[i][:], " ")))
		s.Write([]byte{'\n'})
	}
	s.Write([]byte(strings.Join(symbols[4][:], " ")))
	if b.won {
		s.Write([]byte{'\n', 'W', 'i', 'n', 'n', 'n', 'e', 'r', '!'})
	}
}

type square struct {
	row, col int
}

func newBoard() *board {
	return &board{
		values: make(map[int]*square, 25),
	}
}

// stamp the given number on this board, and check to see if the board won.
func (b *board) stamp(n int) {
	sq := b.values[n]
	if sq == nil {
		return
	}

	delete(b.values, n)
	b.stamped[sq.row][sq.col] = true

	fullColumn := true
	for r := 0; r < 5; r++ {
		if !b.stamped[r][sq.col] {
			fullColumn = false
			break
		}
	}

	fullRow := true
	for c := 0; c < 5; c++ {
		if !b.stamped[sq.row][c] {
			fullRow = false
			break
		}
	}

	if fullColumn || fullRow {
		b.won = true
	}
}

// sum returns the total of all unstamped squares on the board.
func (b *board) sum() int {
	sum := 0
	for k := range b.values {
		sum += k
	}
	return sum
}
