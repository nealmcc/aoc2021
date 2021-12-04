package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	turns, boards, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	i, n := play(boards, turns)
	fmt.Printf("part1: board %d wins on turn %d with score %d\n",
		i, n, boards[i].sum()*turns[n])
}

func read(r io.Reader) (turns []int, boards []*board, err error) {
	s := bufio.NewScanner(r)

	turns, err = readTurns(s)
	if err != nil {
		return
	}

	// discard blank line
	s.Scan()

	boards, err = readBoards(s)
	return
}

func readTurns(s *bufio.Scanner) ([]int, error) {
	if !s.Scan() {
		return nil, errors.New("input too short")
	}
	parts := strings.Split(s.Text(), ",")
	turns := make([]int, 0, 16)
	for _, a := range parts {
		i, err := strconv.Atoi(a)
		if err != nil {
			return nil, err
		}
		turns = append(turns, i)
	}
	return turns, nil
}

func readBoards(s *bufio.Scanner) ([]*board, error) {
	boards := make([]*board, 0, 8)

	for s.Scan() {
		b := newBoard()
		for r := 0; r < 5; r++ {
			c := 0
			text := strings.Split(s.Text(), " ")
			for _, a := range text {
				if len(a) == 0 {
					continue
				}
				n, err := strconv.Atoi(a)
				if err != nil {
					return nil, err
				}
				b.values[n] = &square{row: r, col: c}
				c++
			}
			// discard blank line
			s.Scan()
		}

		boards = append(boards, b)
	}

	return boards, nil
}

// play the given turns for all the boards, until one board wins.
// return the index of the winning board, and the turn number starting from 0
func play(boards []*board, turns []int) (winner, turn int) {
	for i1, n := range turns {
		fmt.Printf("\nround %2d: %2d\n", i1, n)
		for i2, b := range boards {
			b.stamp(n)
			fmt.Println(b)
			fmt.Println()
			if b.won {
				return i2, i1
			}
		}
	}
	return -1, -1
}
