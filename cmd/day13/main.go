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

	paper, folds, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	paper.fold(folds[0])
	fmt.Println("part1", len(paper))

	for _, fn := range folds[1:] {
		paper.fold(fn)
	}
	fmt.Println("part2")
	fmt.Printf("%v", paper)
}

func read(r io.Reader) (paper, []foldFunc, error) {
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

func parseGrid(s *bufio.Scanner) (paper, error) {
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

func parseFolds(s *bufio.Scanner) ([]foldFunc, error) {
	folds := make([]foldFunc, 0, 8)
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

type paper map[v.Coord]bool

// compile-time interface check
var _ fmt.Formatter = paper{}

func (p *paper) fold(fn foldFunc) {
	for dot := range *p {
		delete(*p, dot)
		dot := fn(dot)
		(*p)[dot] = true
	}
}

func (p paper) Format(f fmt.State, verb rune) {
	var max v.Coord
	for dot := range p {
		if dot.X > max.X {
			max.X = dot.X
		}
		if dot.Y > max.Y {
			max.Y = dot.Y
		}
	}

	width, height := max.X+1, max.Y+1

	code := make([][]byte, height)
	for y := 0; y < height; y++ {
		code[y] = []byte(strings.Repeat(" ", width))
	}

	for dot := range p {
		code[dot.Y][dot.X] = '#'
	}

	for _, row := range code {
		f.Write(row)
		f.Write([]byte{'\n'})
	}
}

type foldFunc func(v.Coord) v.Coord

func foldLeft(x int) foldFunc {
	return func(pos v.Coord) v.Coord {
		if pos.X <= x {
			return pos
		}
		pos.X = 2*x - pos.X
		return pos
	}
}

func foldUp(y int) foldFunc {
	return func(pos v.Coord) v.Coord {
		if pos.Y <= y {
			return pos
		}
		pos.Y = 2*y - pos.Y
		return pos
	}
}
