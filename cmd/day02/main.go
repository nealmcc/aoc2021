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

	"github.com/nealmcc/aoc2021/pkg/vector"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	start := vector.Coord{}
	coord, err := part1(start, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 1 product", coord.X*coord.Y)
}

func part1(start vector.Coord, in io.Reader) (vector.Coord, error) {
	s := bufio.NewScanner(in)

	curr := start
	for s.Scan() {
		delta, err := parseCommand(s.Text())
		if err != nil {
			return vector.Coord{}, err
		}
		curr = vector.Add(curr, delta)
	}

	if err := s.Err(); err != nil {
		return vector.Coord{}, err
	}

	return curr, nil
}

// parseCommand is used to interpret a string like 'up 3' and translate
// that into a vector like { X: 0, Y:-3 }.
// Note that for this coordinate system, down is positive Y.
func parseCommand(cmd string) (vector.Coord, error) {
	parts := strings.Split(cmd, " ")

	var v vector.Coord

	switch parts[0] {
	case "forward":
		v.X = 1
	case "up":
		v.Y = -1
	case "down":
		v.Y = 1
	default:
		return v, errors.New("invalid command")
	}

	sc, err := strconv.Atoi(parts[1])
	if err != nil {
		return v, err
	}
	v = vector.Scale(v, sc)
	return v, nil
}
