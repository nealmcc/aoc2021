package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	nav, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1, next := part1(nav)
	fmt.Println("part1", p1)

	p2 := part2(next)
	fmt.Println("part2", p2)
}

func read(r io.Reader) ([]line, error) {
	s := bufio.NewScanner(r)

	lines := make([]line, 0, 16)
	for s.Scan() {
		lines = append(lines, line{raw: []byte(s.Text())})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// part1 returns the answer for part 1, and filters out the corrupt lines
func part1(nav []line) (int, []line) {
	next := make([]line, 0, 16)
	sum := 0
	for _, l := range nav {
		err := l.parse()
		if err, ok := err.(ErrCorrupted); ok && err.IsCorrupt() {
			sum += err.Score()
		} else {
			next = append(next, l)
		}
	}
	return sum, next
}

// part2 returns the answer for part 2, and assumes that all incoming
// lines are incomplete
func part2(nav []line) int {
	scores := make([]int, 0, 32)
	for _, l := range nav {
		rest := l.suggest()
		sc := value(rest)
		scores = append(scores, sc)
	}
	sort.Ints(scores)
	mid := len(scores) / 2
	return scores[mid]
}
