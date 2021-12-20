package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/fish"
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
func read(r io.Reader) ([]fish.Node, error) {
	s := bufio.NewScanner(r)

	nodes := make([]fish.Node, 0, 8)
	for s.Scan() {
		n, err := fish.New(s.Text())
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

// part1 repeats a process of adding nodes, and reducing the sum before adding the next.
func part1(nodes []fish.Node) (fish.Node, error) {
	if len(nodes) == 0 {
		return nil, errors.New("empty input")
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	s := fish.Stack{}
	s.Push(nodes[0])

	for i := 1; i < len(nodes); i++ {
		top := s.Pop()
		pair := &fish.Add{
			L: top,
			R: nodes[i],
		}
		fish.Reduce(pair)
		s.Push(pair)
	}
	return s.Pop(), nil
}
