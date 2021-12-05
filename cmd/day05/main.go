package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/nealmcc/aoc2021/pkg/vector"
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

	// diag, err := render(segments)
	// p1 := part1(diag)

	// fmt.Printf("%v\npart1: %d\n", diag, p1)
}

func read(r io.Reader) ([]vector.Segment, error) {
	s := bufio.NewScanner(r)

	segments := make([]vector.Segment, 0, 16)
	for s.Scan() {
		coords, err := vector.ParseCoords(strings.Split(s.Text(), " -> ")...)
		if err != nil {
			return nil, err
		}
		segments = append(segments, vector.Segment{A: coords[0], B: coords[1]})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}
