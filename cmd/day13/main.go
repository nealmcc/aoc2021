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

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	dots, folds, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(dots, folds[0])
	fmt.Println("part1", p1)
}

func read(r io.Reader) (map[v.Coord]bool, []fold, error) {
	s := bufio.NewScanner(r)

	dots, err := parseGrid(s)
	if err != nil {
		return nil, nil, err
	}

	folds, err := parseFolds(s)
	if err != nil {
		return nil, nil, err
	}

	return dots, folds, nil
}

func parseGrid(s *bufio.Scanner) (map[v.Coord]bool, error) {
	grid := make(map[v.Coord]bool)
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			break
		}
		pos, err := v.ParseCoord(line)
		if err != nil {
			return nil, err
		}
		grid[pos] = true
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return grid, nil
}

func parseFolds(s *bufio.Scanner) ([]fold, error) {
	folds := make([]fold, 0, 8)
	for s.Scan() {
		line := strings.TrimPrefix(s.Text(), "fold along ")
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, errors.New("malformed fold instruction")
		}
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		switch parts[0] {
		case "x":
			folds = append(folds, foldLeft(value))
		case "y":
			folds = append(folds, foldUp(value))
		default:
			return nil, errors.New("malformed fold instruction")
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return folds, nil
}

// part1 solves part 1 of the puzzle
func part1(dots map[v.Coord]bool, fn fold) int {
	for dot := range dots {
		delete(dots, dot)
		dot := fn(dot)
		dots[dot] = true
	}
	return len(dots)
}

type fold func(v.Coord) v.Coord

func foldLeft(x int) fold {
	return func(pos v.Coord) v.Coord {
		if pos.X <= x {
			return pos
		}
		pos.X = 2*x - pos.X
		return pos
	}
}

func foldUp(y int) fold {
	return func(pos v.Coord) v.Coord {
		if pos.Y <= y {
			return pos
		}
		pos.Y = 2*y - pos.Y
		return pos
	}
}
