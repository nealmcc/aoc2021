package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/nealmcc/aoc2021/pkg/vector/stack"
)

const input = `8271653836
7567626775
2315713316
6542655315
2453637333
1247264328
2325146614
2115843171
6182376282
2384738675
`

// main solves part 1 and part 2, reading from the above input
func main() {
	in := strings.NewReader(input)

	octopii, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(&octopii)
	fmt.Println("part1", p1)

	in.Seek(0, io.SeekStart)
	octopii, err = read(in)
	if err != nil {
		log.Fatal(err)
	}

	p2 := part2(&octopii)
	fmt.Println("part2", p2)
}

type grid map[v.Coord]byte

func read(r io.Reader) (grid, error) {
	s := bufio.NewScanner(r)

	g := make(grid)
	y := 0
	for ; s.Scan(); y++ {
		for x, b := range s.Bytes() {
			g[v.Coord{X: x, Y: y}] = b - '0'
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return g, nil
}

// part1 determines how many octopii in total have flashed after 100 steps
func part1(octopii *grid) int {
	sum := 0

	for i := 0; i < 100; i++ {
		sum += octopii.step()
	}

	return sum
}

// part2 determines the first turn that all octopii will flash on the same step
func part2(octopii *grid) int {
	n := 1
	for octopii.step() != 100 {
		n++
	}
	return n
}

// step performs one step of octopus energy adjustment, returning the number
// of octopii that flashed during that step.
func (g *grid) step() int {
	ready := &stack.Coord{}
	for oct := range *g {
		(*g)[oct]++
		if (*g)[oct] == 10 {
			ready.Push(oct)
		}
	}

	done := make(map[v.Coord]bool, 8)
	for ready.Length() > 0 {
		curr := ready.Pop()
		next := g.flash(curr)
		done[curr] = true
		for _, oct := range next {
			if !done[oct] {
				ready.Push(oct)
			}
		}
	}

	for oct := range done {
		(*g)[oct] = 0
	}
	return len(done)
}

// flash causes an octopus at the give position to flash, and returns
// a slice of its neighbours which should flash next.
func (g *grid) flash(pos v.Coord) []v.Coord {
	next := make([]v.Coord, 0, 8)
	for _, c := range neighbours(pos) {
		if _, ok := (*g)[c]; !ok {
			continue
		}

		(*g)[c]++
		if (*g)[c] == 10 {
			next = append(next, c)
		}
	}
	return next
}

// neighbours returns the 8 adjacent coordinates to p
func neighbours(p v.Coord) []v.Coord {
	return []v.Coord{
		{X: p.X - 1, Y: p.Y - 1},
		{X: p.X, Y: p.Y - 1},
		{X: p.X + 1, Y: p.Y - 1},

		{X: p.X - 1, Y: p.Y},
		{X: p.X + 1, Y: p.Y},

		{X: p.X - 1, Y: p.Y + 1},
		{X: p.X, Y: p.Y + 1},
		{X: p.X + 1, Y: p.Y + 1},
	}
}
