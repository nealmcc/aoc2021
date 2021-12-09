package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	land, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(land)
	fmt.Println("part1", p1)
}

func read(r io.Reader) (terrain, error) {
	s := bufio.NewScanner(r)

	t := make(map[v.Coord]int)
	y := 0
	for ; s.Scan(); y++ {
		for x, b := range s.Bytes() {
			t[v.Coord{X: x, Y: y}] = int(b - '0')
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return terrain(t), nil
}

type terrain map[v.Coord]int

func (t terrain) isLow(pos v.Coord) bool {
	center, ok := t[pos]
	if !ok {
		return false
	}

	for x := pos.X - 1; x <= pos.X+1; x++ {
		for y := pos.Y - 1; y <= pos.Y+1; y++ {
			n, ok := t[v.Coord{X: x, Y: y}]
			if !ok {
				continue
			}
			if n < center {
				return false
			}
		}
	}
	return true
}

func part1(t terrain) int {
	sum := 0

	for k, v := range t {
		if t.isLow(k) {
			sum += v + 1
		}
	}

	return sum
}
