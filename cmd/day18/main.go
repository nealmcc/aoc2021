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

	_, err = read(in)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Println("part1:", 0)

	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning a slice of fish trees. Each line of input
// produces one element in the slice.
func read(r io.Reader) ([]ast.Number, error) {
	s := bufio.NewScanner(r)

	numbers := make([]ast.Number, 0, 8)
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
