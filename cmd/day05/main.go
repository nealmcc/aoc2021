package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	segments, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	d1 := render(segments, false)
	fmt.Printf("part1: %d\n", count(d1))

	d2 := render(segments, true)
	fmt.Printf("part2: %d\n", count(d2))
}

// segment is a line segment from a to b.
type segment struct {
	a, b v.Coord
}

// bitmap stores how many fissures are at each coordinate
type bitmap struct {
	points map[v.Coord]int
}

// read a list of line segments from the given input
func read(r io.Reader) ([]segment, error) {
	s := bufio.NewScanner(r)

	segments := make([]segment, 0, 16)
	for s.Scan() {
		coords, err := v.ParseCoords(strings.Split(s.Text(), " -> ")...)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment{a: coords[0], b: coords[1]})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}

// render plots the list of segments, creating a bitmap showing how many
// fissures are present at each position.  The includeDiag parameteter
// determines whether to include diagonal fissures (part 2) or not (part 1)
func render(segments []segment, includeDiag bool) bitmap {
	points := make(map[v.Coord]int)

	for _, seg := range segments {
		delta := v.Sub(seg.b, seg.a)
		if !includeDiag && delta.X != 0 && delta.Y != 0 {
			continue
		}

		curr := seg.a
		points[curr]++

		unit, _ := v.Reduce(delta)
		for curr != seg.b {
			curr = v.Add(curr, unit)
			points[curr]++
		}
	}

	return bitmap{
		points: points,
	}
}

// count the number of points in the given bitmap that have 2 or more fissures
func count(b bitmap) int {
	sum := 0
	for _, n := range b.points {
		if n >= 2 {
			sum++
		}
	}
	return sum
}
