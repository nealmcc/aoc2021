package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/fishtree"
)

// main solves the both part 1 and part 2, reading from input.txt
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
func read(r io.Reader) ([]fishtree.Node, error) {
	s := bufio.NewScanner(r)

	nodes := make([]fishtree.Node, 0, 8)
	for s.Scan() {
		n, err := fishtree.New(s.Text())
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
