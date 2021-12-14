package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	polymer, rules, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(polymer, rules)
	fmt.Println("part1", p1)
}

// read assumes that the rules lines are all well-formed
func read(r io.Reader) (*compound, []rule, error) {
	s := bufio.NewScanner(r)

	if !s.Scan() {
		return nil, nil, errors.New("empty input")
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	polymer, err := newCompound(s.Text())
	if err != nil {
		return nil, nil, err
	}

	// skip blank line
	s.Scan()

	rules := make([]rule, 0, 16)
	for s.Scan() {
		row := s.Bytes()
		rules = append(rules, rule{
			left:  row[0],
			right: row[1],
			mid:   row[6],
		})
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	return polymer, rules, nil
}

type compound struct {
	pairs    map[pair]int
	elements map[byte]int
}

func newCompound(in string) (*compound, error) {
	if len(in) < 2 {
		return nil, errors.New("a compound must have at least two elements")
	}

	c := compound{
		pairs:    make(map[pair]int, 4),
		elements: make(map[byte]int, 4),
	}
	curr := in[0]
	c.elements[curr] = 1
	for i := 1; i < len(in); i++ {
		next := in[i]
		c.pairs[pair{left: curr, right: next}] += 1
		curr = next
		c.elements[curr] += 1
	}

	return &c, nil
}

type rule struct {
	left, right byte
	mid         byte
}

type pair struct {
	left, right byte
}

// part1 replicates the template 10 times, following the given rules, and then
// returns the difference between the quantity of hte most common element,
// and the least common element in the resulting polymer.
func part1(c *compound, rules []rule) int {
	for i := 0; i < 10; i++ {
		c.replicate(rules)
	}

	anykey := getAnyKey(c.elements)
	min, max := c.elements[anykey], c.elements[anykey]

	for _, count := range c.elements {
		if count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}
	return max - min
}

// getAnyKey returns any key from the map - we don't care which one.
func getAnyKey(m map[byte]int) byte {
	for k := range m {
		return k
	}
	return 0
}

func (c *compound) replicate(rules []rule) {
	type op struct {
		p     pair
		count int
	}

	removals := make([]op, 0, 16)
	insertions := make([]op, 0, 16)

	for _, r := range rules {
		remove := pair{left: r.left, right: r.right}
		count := c.pairs[remove]
		removals = append(removals, op{p: remove, count: count})

		add1 := pair{left: r.left, right: r.mid}
		add2 := pair{left: r.mid, right: r.right}
		insertions = append(insertions, op{p: add1, count: count})
		insertions = append(insertions, op{p: add2, count: count})

		c.elements[r.mid] += count
	}

	for _, r := range removals {
		c.pairs[r.p] -= r.count
	}

	for _, r := range insertions {
		c.pairs[r.p] += r.count
	}
}
