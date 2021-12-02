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

	sub1 := submarine{}
	err = part1(&sub1, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 1 product", sub1.coord.X*sub1.coord.Y)

	in.Seek(0, io.SeekStart)

	sub2 := submarine{}
	part2(&sub2, in)
	fmt.Println("part 2 product", sub2.coord.X*sub2.coord.Y)
}

func part1(sub *submarine, in io.Reader) error {
	s := bufio.NewScanner(in)

	for s.Scan() {
		if err := sub.movePart1(s.Text()); err != nil {
			return err
		}
	}

	return s.Err()
}

func part2(sub *submarine, in io.Reader) error {
	s := bufio.NewScanner(in)
	for s.Scan() {
		if err := sub.movePart2(s.Text()); err != nil {
			return err
		}
	}

	return s.Err()
}

// submarine represents a submarine with a position and (for part2) an aim.
type submarine struct {
	// coord is the submarine's position. Down is positive Y.
	coord vector.Coord
	// aim is used for part2, when adjusting the submarine's up/down aim.
	aim int
}

// movePart1 parses and executes the given command, following the rules for part1.
func (s *submarine) movePart1(cmd string) error {
	parts := strings.Split(cmd, " ")

	dir := parts[0]
	scale, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	switch dir {
	case "forward":
		s.coord.X += scale
	case "up":
		s.coord.Y -= scale
	case "down":
		s.coord.Y += scale
	default:
		return errors.New("invalid command")
	}

	return nil
}

// movePart2 parses and executes the given command, following the rules for part2.
func (s *submarine) movePart2(cmd string) error {
	parts := strings.Split(cmd, " ")

	dir := parts[0]
	scale, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	switch dir {
	case "forward":
		s.coord.X += scale
		s.coord.Y += s.aim * scale
	case "up":
		s.aim -= scale
	case "down":
		s.aim += scale
	default:
		return errors.New("invalid command")
	}

	return nil
}
