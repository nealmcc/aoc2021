package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	caves, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(caves)
	fmt.Println("part1", p1)
}

type node struct {
	id  string
	out map[string]*node
}

func read(r io.Reader) (map[string]*node, error) {
	s := bufio.NewScanner(r)

	caves := make(map[string]*node, 24)

	getOrCreate := func(id string) *node {
		n, ok := caves[id]
		if !ok {
			n = &node{id: id, out: make(map[string]*node, 2)}
			caves[id] = n
		}
		return n
	}

	for s.Scan() {
		parts := strings.Split(s.Text(), "-")
		left, right := getOrCreate(parts[0]), getOrCreate(parts[1])
		left.out[right.id] = right
		right.out[left.id] = left
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return caves, nil
}

// part1 solves part1 of the puzzle
func part1(caves map[string]*node) int {
	rules := &p1Rules{visited: make(map[string]bool)}
	paths := pathsFrom(caves, caves["start"], caves["end"], rules)
	return len(paths)
}

// pathsFrom lists all the possible paths through the cave system from
// start to end, using the given cave rules to determine which caves
// are allowed to be visited.
func pathsFrom(caves map[string]*node, start, end *node, s caveRules) []string {
	if start.id == end.id {
		return []string{end.id}
	}

	s.visit(start.id)
	paths := make([]string, 0, 2)
	for key, child := range start.out {
		if s.canVisit(key) {
			childPaths := pathsFrom(caves, child, end, s.clone())
			for _, p := range childPaths {
				paths = append(paths, start.id+","+p)
			}
		}
	}

	return paths
}

// caveRules is the interface that determines the rules for visiting caves
type caveRules interface {
	canVisit(id string) bool
	visit(id string)
	clone() caveRules
}

// p1Rules defines the part 1 rules
type p1Rules struct {
	visited map[string]bool
}

// compile-time interface check
var _ caveRules = new(p1Rules)

func (p1 *p1Rules) canVisit(id string) bool {
	if id[0] < 'a' {
		return true
	}
	return !p1.visited[id]
}

func (p1 *p1Rules) visit(id string) {
	p1.visited[id] = true
}

func (p1 *p1Rules) clone() caveRules {
	visited := make(map[string]bool, len(p1.visited))
	for key, val := range p1.visited {
		visited[key] = val
	}
	return &p1Rules{visited}
}