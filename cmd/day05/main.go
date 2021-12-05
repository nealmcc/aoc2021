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
	fmt.Println(segments)

	d1 := render(segments, false)
	fmt.Printf("part1: %d\n", count(d1))

	d2 := render(segments, true)
	fmt.Printf("part2: %d\n", count(d2))
}

func read(r io.Reader) ([]v.Segment, error) {
	s := bufio.NewScanner(r)

	segments := make([]v.Segment, 0, 16)
	for s.Scan() {
		coords, err := v.ParseCoords(strings.Split(s.Text(), " -> ")...)
		if err != nil {
			return nil, err
		}
		segments = append(segments, v.Segment{A: coords[0], B: coords[1]})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}

type bitmap struct {
	// points stores how many fissures overlap the given coordinate
	points map[v.Coord]int
}

func render(segments []v.Segment, includeDiag bool) bitmap {
	points := make(map[v.Coord]int)

	for _, seg := range segments {
		delta := v.Sub(seg.B, seg.A)
		if !includeDiag && delta.X != 0 && delta.Y != 0 {
			continue
		}

		curr := seg.A
		points[curr] = points[curr] + 1

		unit, _ := v.Reduce(delta)
		for curr != seg.B {
			curr = v.Add(curr, unit)
			points[curr] = points[curr] + 1
		}
	}

	return bitmap{
		points: points,
	}
}

func count(b bitmap) int {
	sum := 0
	for _, n := range b.points {
		if n >= 2 {
			sum++
		}
	}
	return sum
}
