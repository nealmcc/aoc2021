package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	polymer, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		polymer.replicate()
	}
	p1 := polymer.magicNumber()

	for i := 0; i < 30; i++ {
		polymer.replicate()
	}
	p2 := polymer.magicNumber()
	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read assumes that the rules lines are all well-formed
func read(r io.Reader) (*compound, error) {
	s := bufio.NewScanner(r)

	if !s.Scan() {
		return nil, errors.New("empty input")
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	polymer, err := newCompound(s.Text())
	if err != nil {
		return nil, err
	}

	// skip blank line
	s.Scan()

	polymer.rules = make([]rule, 0, 100)
	for s.Scan() {
		row := s.Bytes()
		polymer.rules = append(polymer.rules, rule{
			left:  row[0] - 'A',
			right: row[1] - 'A',
			mid:   row[6] - 'A',
		})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return polymer, nil
}

type compound struct {
	pairs    [26][26]int
	elements [26]int
	rules    []rule
}

// rule defines the element that will be inserted between a pair of elements.
type rule struct {
	left, right byte
	mid         byte
}

// newCompound creates a new polymer using the given starting elements.
// The compound will not have any rules - those need to be added separately.
func newCompound(in string) (*compound, error) {
	if len(in) < 2 {
		return nil, errors.New("a compound must have at least two elements")
	}

	c := compound{
		pairs:    [26][26]int{},
		elements: [26]int{},
	}
	curr := in[0] - 'A'
	c.elements[curr] = 1
	for i := 1; i < len(in); i++ {
		next := in[i] - 'A'
		c.pairs[curr][next] += 1
		curr = next
		c.elements[curr] += 1
	}

	return &c, nil
}

// replicate performs one step of growth for the polymer
func (c *compound) replicate() {
	type op struct {
		l, r byte
		n    int
	}

	remove := make([]op, 0, 100)
	add := make([]op, 0, 200)

	for _, r := range c.rules {
		count := c.pairs[r.left][r.right]
		remove = append(remove, op{l: r.left, r: r.right, n: count})
		add = append(add, op{l: r.left, r: r.mid, n: count})
		add = append(add, op{l: r.mid, r: r.right, n: count})

		c.elements[r.mid] += count
	}

	for _, op := range remove {
		c.pairs[op.l][op.r] -= op.n
	}

	for _, op := range add {
		c.pairs[op.l][op.r] += op.n
	}
}

// magicNumber is the difference between the quantity of the most common
// element and the least common element in the compound.
func (c *compound) magicNumber() int {
	min, max := 1<<63-1, 0

	for _, count := range c.elements {
		if count > 0 && count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}
	return max - min
}
