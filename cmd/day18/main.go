package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/ast"
)

// main solves both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	lines, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1, err := part1(lines)
	if err != nil {
		log.Fatal(err)
	}

	p2, err := part2(lines)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning a slice of fish trees. Each line of input
// produces one element in the slice.
func read(r io.Reader) ([]*ast.Tree, error) {
	s := bufio.NewScanner(r)

	numbers := make([]*ast.Tree, 0, 8)
	for s.Scan() {
		n, err := ast.New(s.Text())
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, n)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return numbers, nil
}

func part1(numbers []*ast.Tree) (int, error) {
	sum, err := ast.Sum(numbers...)
	if err != nil {
		return 0, err
	}
	return sum.Magnitude(), nil
}

func part2(numbers []*ast.Tree) (int, error) {
	var best int

	for i := 0; i < len(numbers); i++ {
		for j := 0; j < len(numbers); j++ {
			if j == i {
				continue
			}
			sum, err := ast.Add(*numbers[i], *numbers[j])
			if err != nil {
				return 0, err
			}
			mag := sum.Magnitude()
			if mag > best {
				best = mag
			}
		}
	}

	return best, nil
}
