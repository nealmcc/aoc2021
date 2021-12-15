package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	cave, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	var (
		from = vector.Coord{X: 0, Y: 0}
		to   = vector.Coord{X: 99, Y: 99}
	)
	p1 := shortestPath(cave, from, to)

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read a cave from the given input
func read(r io.Reader) (cavern, error) {
	s := bufio.NewScanner(r)

	if !s.Scan() {
		return nil, errors.New("empty input")
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	line1 := s.Bytes()
	cave := newCavern(len(line1))
	copy(cave[0], line1)

	i := 1
	for s.Scan() {
		copy(cave[i], s.Bytes())
		i++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return cave, nil
}

// cavern is a square space with coordinates ranging from (0, 0) at the top left
// to (maxCol, maxRow) at the bottom right.
// Each position in the cave has a variable risk level from 0 to 9.
type cavern [][]byte

// newCavern creates a new n x n cavern with 0 risk level at each position.
func newCavern(n int) cavern {
	c := cavern(make([][]byte, n))

	for i := 0; i < n; i++ {
		c[i] = make([]byte, n)
	}

	return c
}

func shortestPath(c cavern, from, to vector.Coord) int {
	return 0
}
