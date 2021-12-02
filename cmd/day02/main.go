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
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := vector.Coord{}
	coord, err := part1(start, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 1 product", coord.X*coord.Y)

	in.Seek(0, io.SeekStart)

	sub := submarine{}
	part2(&sub, in)
	fmt.Println("part 2 product", sub.coord.X*sub.coord.Y)
}

func part1(start vector.Coord, in io.Reader) (vector.Coord, error) {
	s := bufio.NewScanner(in)

	// parsePart1 converts 'up 3' into a vector { X: 0, Y:-3 }.
	// Note that for this coordinate system, down is positive Y.
	parsePart1 := func(cmd string) (vector.Coord, error) {
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

	curr := start
	for s.Scan() {
		delta, err := parsePart1(s.Text())
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

func part2(sub *submarine, in io.Reader) error {
	s := bufio.NewScanner(in)

	for s.Scan() {
		if err := sub.do(s.Text()); err != nil {
			return err
		}
	}

	if err := s.Err(); err != nil {
		return err
	}
	return nil
}

type submarine struct {
	coord vector.Coord
	aim   int
}

// do parses and executes the given command, following the rules for part2.
func (s *submarine) do(cmd string) error {
	parts := strings.Split(cmd, " ")
	dir := parts[0]
	scale, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	switch dir {
	case "down":
		s.aim += scale
	case "up":
		s.aim -= scale
	case "forward":
		s.coord.X += scale
		s.coord.Y += s.aim * scale
	default:
		return errors.New("invalid command")
	}
	return nil
}
