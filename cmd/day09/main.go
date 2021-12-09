package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

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

	p2 := part2(land)
	fmt.Println("part2", p2)
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

func part1(t terrain) int {
	sum := 0
	for _, pos := range t.lowPoints() {
		sum += t[pos] + 1
	}
	return sum
}

func part2(t terrain) int {
	basins := t.basins()
	sort.Slice(basins, func(i, j int) bool {
		return len(basins[i]) > len(basins[j])
	})

	prod := 1
	for _, b := range basins[:3] {
		prod *= len(b)
	}

	return prod
}

// lowPoints returns a slice of all the lowPoints of this terrain
func (t terrain) lowPoints() []v.Coord {
	low := make([]v.Coord, 0, 16)
	for k := range t {
		if t.isLow(k) {
			low = append(low, k)
		}
	}
	return low
}

// isLow determines if the given position is a low point for this terrain
func (t terrain) isLow(pos v.Coord) bool {
	center, ok := t[pos]
	if !ok {
		return false
	}

	for _, p := range neighbours(pos) {
		n, ok := t[p]
		if !ok {
			continue
		}
		if n < center {
			return false
		}
	}
	return true
}

// neighbours returns the 4 adjacent coordinates to p
func neighbours(p v.Coord) []v.Coord {
	return []v.Coord{
		{X: p.X - 1, Y: p.Y},
		{X: p.X + 1, Y: p.Y},
		{X: p.X, Y: p.Y - 1},
		{X: p.X, Y: p.Y + 1},
	}
}

// basins finds all of the contiguous basins of this terrain
func (t terrain) basins() []terrain {
	basins := make([]terrain, 0, 8)

	for _, p := range t.lowPoints() {
		delete(t, p)
		b := newBasin(p, t[p])
		basins = append(basins, terrain(b))
	}

	for pos, h := range t {
		if h == 9 {
			delete(t, pos)
			continue
		}
		for _, b := range basins {
			if b.connectsTo(pos) {
				b[pos] = h
				delete(t, pos)
				continue
			}
		}
	}

	for len(t) > 0 {
		for pos, h := range t {
			for _, b := range basins {
				if b.connectsTo(pos) {
					b[pos] = h
					delete(t, pos)
					continue
				}
			}
		}
	}

	return basins
}

// newBasin creates a new terrain with the given position and height
func newBasin(pos v.Coord, height int) terrain {
	b := make(map[v.Coord]int)
	b[pos] = height
	return terrain(b)
}

// connectsTo returns true if the given position is directly adjacent to this terrain
func (t terrain) connectsTo(pos v.Coord) bool {
	for _, p := range neighbours(pos) {
		if _, ok := t[p]; ok {
			return true
		}
	}
	return false
}
