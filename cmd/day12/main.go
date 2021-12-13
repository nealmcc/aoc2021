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

	p2 := part2(caves)
	fmt.Println("part2", p2)
}

type network struct {
	names []string
	nodes []*node
}

type node struct {
	id  int
	out []*node
}

func (net *network) getOrCreate(name string) *node {
	for k, v := range net.names {
		if v == name {
			return net.nodes[k]
		}
	}
	cave := &node{
		id:  len(net.nodes),
		out: []*node{},
	}
	net.names = append(net.names, name)
	net.nodes = append(net.nodes, cave)

	return cave
}

func read(r io.Reader) (*network, error) {
	s := bufio.NewScanner(r)

	net := network{
		names: make([]string, 0, 16),
		nodes: make([]*node, 0, 16),
	}

	for s.Scan() {
		parts := strings.Split(s.Text(), "-")
		left, right := net.getOrCreate(parts[0]), net.getOrCreate(parts[1])
		left.out = append(left.out, right)
		right.out = append(right.out, left)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return &net, nil
}

// part1 solves part 1 of the puzzle
func part1(caves *network) int {
	rules := newP1(len(caves.names))
	start := caves.getOrCreate("start")
	end := caves.getOrCreate("end")
	return countPaths(caves, start, end, rules)
}

// part2 solves part 2 of the puzzle
func part2(caves *network) int {
	rules := newP2(len(caves.names))
	start := caves.getOrCreate("start")
	end := caves.getOrCreate("end")
	return countPaths(caves, start, end, rules)
}

// countPaths counts the distinct paths through the cave system from
// start to end, using the given cave rules to determine which caves
// are allowed to be visited.
func countPaths(caves *network, start, end *node, rules caveRules) int {
	stack := Stack{}
	stack.Push(state{cave: start, rules: rules})

	var sum int
	for stack.Length() > 0 {
		x := stack.Pop()
		cave, rules := x.cave, x.rules
		rules.visit(cave.id, caves.names[cave.id])
		for _, child := range cave.out {
			if child.id == end.id {
				sum++
				continue
			}
			if rules.canVisit(child.id, caves.names[child.id]) {
				stack.Push(state{
					cave:  child,
					rules: rules.clone(),
				})
			}
		}
	}
	return sum
}

// caveRules is the interface that determines the rules for visiting caves
type caveRules interface {
	canVisit(id int, name string) bool
	visit(id int, name string)
	clone() caveRules
}

// p1Rules defines the part 1 rules
type p1Rules []bool

func newP1(numCaves int) *p1Rules {
	rules := make(p1Rules, numCaves)
	return &rules
}

// compile-time interface check
var _ caveRules = new(p1Rules)

func (p1 *p1Rules) canVisit(id int, name string) bool {
	if name[0] < 'a' {
		return true
	}
	return !(*p1)[id]
}

func (p1 *p1Rules) visit(id int, _ string) {
	(*p1)[id] = true
}

func (p1 *p1Rules) clone() caveRules {
	next := make(p1Rules, len(*p1))
	copy(next, *p1)
	return &next
}

// p2Rules defines the part 2 rules
type p2Rules struct {
	did2x   bool
	visited []bool
}

func newP2(numCaves int) *p2Rules {
	return &p2Rules{
		visited: make([]bool, numCaves),
	}
}

// compile-time interface check
var _ caveRules = new(p2Rules)

func (p2 *p2Rules) canVisit(id int, name string) bool {
	switch {
	case name[0] < 'a':
		return true
	case !p2.visited[id]:
		return true
	case name == "start" || name == "end":
		return false
	default:
		return !p2.did2x
	}
}

func (p2 *p2Rules) visit(id int, name string) {
	if p2.visited[id] && name[0] >= 'a' {
		p2.did2x = true
	}
	p2.visited[id] = true
}

func (p2 *p2Rules) clone() caveRules {
	next := p2Rules{
		did2x:   p2.did2x,
		visited: make([]bool, len(p2.visited)),
	}
	copy(next.visited, p2.visited)
	return &next
}
